package io

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
)

func CopyPath(srcPath, dstPath string, f os.FileInfo, copyMode Mode) error {
	stat, ok := f.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("Unable to get raw syscall.Stat_t data for %s", srcPath)
	}

	isHardlink := false

	switch mode := f.Mode(); {
	case mode.IsRegular():
		if copyMode == Hardlink {
			isHardlink = true
			if err := os.Link(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyRegular(srcPath, dstPath, f); err != nil {
				return err
			}
		}

	case mode.IsDir():
		if err := os.Mkdir(dstPath, f.Mode()); err != nil && !os.IsExist(err) {
			return err
		}
	case mode&os.ModeSymlink != 0:
		link, err := os.Readlink(srcPath)
		if err != nil {
			return err
		}

		if err := os.Symlink(link, dstPath); err != nil {
			return err
		}

	case mode&os.ModeNamedPipe != 0:
		fallthrough

	case mode&os.ModeSocket != 0:
		if err := unix.Mkfifo(dstPath, uint32(stat.Mode)); err != nil {
			return err
		}
	case mode&os.ModeDevice != 0:
		if err := unix.Mknod(dstPath, uint32(stat.Mode), int(stat.Rdev)); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown file type (%d / %s) for %s", f.Mode(), f.Mode().String(), srcPath)

	}

	// Everything below is copying metadata from src to dst. All this metadata
	// already shares an inode for hardlinks.
	if isHardlink {
		return nil
	}

	if err := os.Lchown(dstPath, int(stat.Uid), int(stat.Gid)); err != nil {
		return err
	}

	isSymlink := f.Mode()&os.ModeSymlink != 0
	// There is no LChmod, so ignore mode for symlink. Also, this
	// must happen after chown, as that can modify the file mode
	if !isSymlink {
		if err := os.Chmod(dstPath, f.Mode()); err != nil {
			return err
		}
	}

	return nil
}

func doCopyWithFileClone(srcFile, dstFile *os.File) error {
	return unix.IoctlFileClone(int(dstFile.Fd()), int(srcFile.Fd()))
}

func doCopyWithFileRange(srcFile, dstFile *os.File, fileinfo os.FileInfo) error {
	amountLeftToCopy := fileinfo.Size()

	for amountLeftToCopy > 0 {
		n, err := unix.CopyFileRange(int(srcFile.Fd()), nil, int(dstFile.Fd()), nil, int(amountLeftToCopy), 0)
		if err != nil {
			return err
		}

		amountLeftToCopy = amountLeftToCopy - int64(n)
	}

	return nil

}
