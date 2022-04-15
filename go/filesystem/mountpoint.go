package filesystem

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	os_ "github.com/kaydxh/golang/go/os"
	filepath_ "github.com/kaydxh/golang/go/path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

var (
	mountsByDevice    map[DeviceNumber]*Mount
	mountsByPath      map[string]*Mount
	mountMutex        sync.Mutex
	mountsInitialized bool
	allMountsByDevice map[DeviceNumber][]*Mount
)

type DeviceNumber uint64

func (num DeviceNumber) String() string {
	return fmt.Sprintf("%d:%d", unix.Major(uint64(num)), unix.Minor(uint64(num)))
}

// getNumberOfContainingDevice returns the device number of the filesystem which
// contains the given file.  If the file is a symlink, it is not dereferenced.
func getNumberOfContainingDevice(path string) (DeviceNumber, error) {
	var stat unix.Stat_t
	if err := unix.Lstat(path, &stat); err != nil {
		return 0, err
	}
	return DeviceNumber(stat.Dev), nil
}

func filesystemLacksMainMountError(deviceNumber DeviceNumber) error {
	return errors.Errorf(
		"Device %q (%v) lacks a \"main\" mountpoint in the current mount namespace, so it's ambiguous where to store the fscrypt metadata.",
		getDeviceName(deviceNumber),
		deviceNumber,
	)
}

type Mount struct {
	Path           string
	FilesystemType string
	Device         string
	DeviceNumber   DeviceNumber
	Subtree        string
	ReadOnly       bool
}

func FindMount(path string) (*Mount, error) {
	mountMutex.Lock()
	defer mountMutex.Unlock()
	if err := loadMountInfo(); err != nil {
		return nil, err
	}
	// First try to find the mount by the number of the containing device.
	deviceNumber, err := getNumberOfContainingDevice(path)
	if err != nil {
		return nil, err
	}
	mnt, ok := mountsByDevice[deviceNumber]

	if ok {
		if mnt == nil {
			mnts, ok := allMountsByDevice[deviceNumber]
			if ok {
				if len(mnts) == 0 {
					return nil, filesystemLacksMainMountError(deviceNumber)
				}

				for _, mnt := range mnts {
					if strings.HasPrefix(path, mnt.Path) {
						return mnt, nil
					}
				}
			}
			return nil, filesystemLacksMainMountError(deviceNumber)
		}

		return mnt, nil
	}
	// The mount couldn't be found by the number of the containing device.
	// Fall back to walking up the directory hierarchy and checking for a
	// mount at each directory path.  This is necessary for btrfs, where
	// files report a different st_dev from the /proc/self/mountinfo entry.
	curPath, err := filepath_.CanonicalizePath(path)
	if err != nil {
		return nil, err
	}
	for {
		mnt := mountsByPath[curPath]
		if mnt != nil {
			return mnt, nil
		}
		// Move to the parent directory unless we have reached the root.
		parent := filepath.Dir(curPath)
		if parent == curPath {
			return nil, errors.Errorf("couldn't find mountpoint containing %q", path)
		}
		curPath = parent
	}
}

// loadMountInfo populates the Mount mappings by parsing /proc/self/mountinfo.
// It returns an error if the Mount mappings cannot be populated.
func loadMountInfo() error {
	if !mountsInitialized {
		file, err := os.Open("/proc/self/mountinfo")
		if err != nil {
			return err
		}
		defer file.Close()
		if err := readMountInfo(file); err != nil {
			return err
		}
		mountsInitialized = true
	}
	return nil
}

// This is separate from loadMountInfo() only for unit testing.
func readMountInfo(r io.Reader) error {
	mountsByDevice = make(map[DeviceNumber]*Mount)
	mountsByPath = make(map[string]*Mount)
	allMountsByDevice = make(map[DeviceNumber][]*Mount)
	allMountsByPath := make(map[string]*Mount)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		mnt := parseMountInfoLine(line)
		if mnt == nil {
			logrus.Warnf("ignoring invalid mountinfo line %q", line)
			continue
		}

		// We can only use mountpoints that are directories for fscrypt.
		isDir, err := os_.IsDir(mnt.Path)
		if err != nil {
			logrus.Errorf("ignoring mountpoint %v because isDir failed", err)
			continue
		}

		if !isDir {
			logrus.Infof("ignoring mountpoint %q because it is not a directory", mnt.Path)
			continue
		}

		// Note this overrides the info if we have seen the mountpoint
		// earlier in the file. This is correct behavior because the
		// mountpoints are listed in mount order.
		allMountsByPath[mnt.Path] = mnt
	}
	// For each filesystem, choose a "main" Mount and discard any additional
	// bind mounts.  fscrypt only cares about the main Mount, since it's
	// where the fscrypt metadata is stored.  Store all the main Mounts in
	// mountsByDevice and mountsByPath so that they can be found later.
	for _, mnt := range allMountsByPath {
		allMountsByDevice[mnt.DeviceNumber] =
			append(allMountsByDevice[mnt.DeviceNumber], mnt)
	}

	for deviceNumber, filesystemMounts := range allMountsByDevice {
		mnt := findMainMount(filesystemMounts)
		mountsByDevice[deviceNumber] = mnt // may store an explicit nil entry
		if mnt != nil {
			mountsByPath[mnt.Path] = mnt
		}
	}
	return nil
}

// For more details, see https://www.kernel.org/doc/Documentation/filesystems/proc.txt
func parseMountInfoLine(line string) *Mount {
	fields := strings.Split(line, " ")
	if len(fields) < 10 {
		return nil
	}

	// Count the optional fields.  In case new fields are appended later,
	// don't simply assume that n == len(fields) - 4.
	n := 6
	for fields[n] != "-" {
		n++
		if n >= len(fields) {
			return nil
		}
	}
	if n+3 >= len(fields) {
		return nil
	}

	var mnt *Mount = &Mount{}
	var err error
	mnt.DeviceNumber, err = newDeviceNumberFromString(fields[2])
	if err != nil {
		return nil
	}
	mnt.Subtree = unescapeString(fields[3])
	mnt.Path = unescapeString(fields[4])
	for _, opt := range strings.Split(fields[5], ",") {
		if opt == "ro" {
			mnt.ReadOnly = true
		}
	}
	mnt.FilesystemType = unescapeString(fields[n+1])
	mnt.Device = getDeviceName(mnt.DeviceNumber)
	return mnt
}

func newDeviceNumberFromString(str string) (DeviceNumber, error) {
	var major, minor uint32
	if count, _ := fmt.Sscanf(str, "%d:%d", &major, &minor); count != 2 {
		return 0, errors.Errorf("invalid device number string %q", str)
	}
	return DeviceNumber(unix.Mkdev(major, minor)), nil
}

// Unescape octal-encoded escape sequences in a string from the mountinfo file.
// The kernel encodes the ' ', '\t', '\n', and '\\' bytes this way.  This
// function exactly inverts what the kernel does, including by preserving
// invalid UTF-8.
func unescapeString(str string) string {
	var sb strings.Builder
	for i := 0; i < len(str); i++ {
		b := str[i]
		if b == '\\' && i+3 < len(str) {
			if parsed, err := strconv.ParseInt(str[i+1:i+4], 8, 8); err == nil {
				b = uint8(parsed)
				i += 3
			}
		}
		sb.WriteByte(b)
	}
	return sb.String()
}

// We get the device name via the device number rather than use the mount source
// field directly.  This is necessary to handle a rootfs that was mounted via
// the kernel command line, since mountinfo always shows /dev/root for that.
// This assumes that the device nodes are in the standard location.
func getDeviceName(num DeviceNumber) string {
	linkPath := fmt.Sprintf("/sys/dev/block/%v", num)
	if target, err := os.Readlink(linkPath); err == nil {
		return fmt.Sprintf("/dev/%s", filepath.Base(target))
	}
	return ""
}

type mountpointTreeNode struct {
	mount    *Mount
	parent   *mountpointTreeNode
	children []*mountpointTreeNode
}

func addUncontainedSubtreesRecursive(dst map[string]bool,
	node *mountpointTreeNode, allUncontainedSubtrees map[string]bool) {
	if allUncontainedSubtrees[node.mount.Subtree] {
		dst[node.mount.Subtree] = true
	}
	for _, child := range node.children {
		addUncontainedSubtreesRecursive(dst, child, allUncontainedSubtrees)
	}
}

func findMainMount(filesystemMounts []*Mount) *Mount {
	// Index this filesystem's mounts by path.  Note: paths are unique here,
	// since non-last mounts were already excluded earlier.
	//
	// Also build the set of all mounted subtrees.
	filesystemMountsByPath := make(map[string]*mountpointTreeNode)
	allSubtrees := make(map[string]bool)
	for _, mnt := range filesystemMounts {
		filesystemMountsByPath[mnt.Path] = &mountpointTreeNode{mount: mnt}
		allSubtrees[mnt.Subtree] = true
	}

	// Divide the mounts into non-overlapping trees of mountpoints.
	for path, mntNode := range filesystemMountsByPath {
		for path != "/" && mntNode.parent == nil {
			path = filepath.Dir(path)
			if parent := filesystemMountsByPath[path]; parent != nil {
				mntNode.parent = parent
				parent.children = append(parent.children, mntNode)
			}
		}
	}

	// Build the set of mounted subtrees that aren't contained in any other
	// mounted subtree.
	allUncontainedSubtrees := make(map[string]bool)
	for subtree := range allSubtrees {
		contained := false
		for t := subtree; t != "/" && !contained; {
			t = filepath.Dir(t)
			contained = allSubtrees[t]
		}
		if !contained {
			allUncontainedSubtrees[subtree] = true
		}
	}

	// Select the root of a mountpoint tree whose mounted subtrees contain
	// *all* mounted subtrees.  Equivalently, select a mountpoint tree in
	// which every uncontained subtree is mounted.
	var mainMount *Mount
	for _, mntNode := range filesystemMountsByPath {
		mnt := mntNode.mount
		if mntNode.parent != nil {
			continue
		}
		uncontainedSubtrees := make(map[string]bool)
		addUncontainedSubtreesRecursive(uncontainedSubtrees, mntNode, allUncontainedSubtrees)
		if len(uncontainedSubtrees) != len(allUncontainedSubtrees) {
			continue
		}
		// If there's more than one eligible mount, they should have the
		// same Subtree.  Otherwise it's ambiguous which one to use.
		if mainMount != nil && mainMount.Subtree != mnt.Subtree {
			logrus.Errorf(
				"Unsupported case: %q (%v) has multiple non-overlapping mounts. This filesystem will be ignored!",
				mnt.Device,
				mnt.DeviceNumber,
			)
			return nil
		}
		// Prefer a read-write mount to a read-only one.
		if mainMount == nil || mainMount.ReadOnly {
			mainMount = mnt
		}
	}
	return mainMount
}
