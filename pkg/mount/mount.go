package mount

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	io_ "github.com/kaydxh/golang/go/io"

	os_ "github.com/kaydxh/golang/go/os"
	exec_ "github.com/kaydxh/golang/go/os/exec"
)

const (
	mountLabelFileName = "mount-label-file-75248"
)

func MountCeph(
	prefixMountPath, userName, address, keyring, workDir string, forceMount bool, timeout int,
) (mountPoint string, err error) {
	keyringFilePath, err := genCephKeyringConfFile(prefixMountPath, userName, keyring)
	if err != nil {
		return "", err
	}

	workDir = filepath.Join("/", workDir)

	cmd := fmt.Sprintf(
		`ceph-fuse -m %s -k %s -n client.%s -o rw.nonempty -r %s %s`,
		address,
		keyringFilePath,
		userName,
		workDir,
		mountPoint,
	)
	return doMountCmd(prefixMountPath, workDir, cmd, forceMount, timeout)
}

func MountCfs(
	prefixMountPath, address, workDir string,
	forceMount bool,
	timeout int,
) (mountPoint string, err error) {

	mountPoint = filepath.Join(prefixMountPath, workDir)
	workDir = filepath.Join("/", workDir)

	cmd := fmt.Sprintf(
		`mount -t nfs -o vers=4 %s:%s %s`,
		address,
		workDir,
		mountPoint,
	)
	return doMountCmd(prefixMountPath, workDir, cmd, forceMount, timeout)
}

func doMountCmd(
	prefixMountPath, workDir, cmd string,
	forceMount bool,
	timeout int,
) (mountPoint string, err error) {

	mountPoint = filepath.Join("/", prefixMountPath, workDir)
	if !forceMount {
		labelFilePath := filepath.Join(mountPoint, mountLabelFileName)
		exist, err := os_.PathExists(labelFilePath)
		if err == nil && exist {
			return mountPoint, nil
		}
	}

	err = os.MkdirAll(mountPoint, 0755)
	if err != nil {
		return "", err
	}

	_, _, err = exec_.Exec(time.Duration(timeout), cmd)
	if err != nil {
		return "", err
	}

	return mountPoint, nil
}

func genCephKeyringConfFile(prefixMountPath, userName, keyring string) (string, error) {
	keyringFileName := userName + ".keyring"
	keyringFilePath := filepath.Join(prefixMountPath, keyringFileName)

	line1 := "[client." + userName + "]"
	line2 := fmt.Sprintf("    key =  %s\n", keyring)
	lines := []string{line1, line2}

	err := io_.WriteFileLines(keyringFilePath, lines, false)
	if err != nil {
		return "", err
	}

	return keyringFilePath, nil
}
