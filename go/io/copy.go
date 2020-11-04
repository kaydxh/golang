package io

import (
	"fmt"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

// Mode indicates whether to use hardlink or copy content
type Mode int

const (
	// Content creates a new file, and copies the content of the file
	Content Mode = iota
	// Hardlink creates a new hardlink to the existing file
	Hardlink
)

func CopyDir(srcDir, dstDir string, copyMode Mode) (err error) {
	return filepath.Walk(srcDir, func(srcPath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Rebase path
		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, relPath)

		return CopyPath(srcPath, dstPath, f, copyMode)
	})

}

func CopyPath(srcPath, dstPath string, f os.FileInfo, copyMode Mode) error {
	stat, ok := f.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("Unable to get raw syscall.Stat_t data for %s", srcPath)
	}

	isHardlink := false

	switch mode := f.Mode(); {
	case mode.IsRegular():

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

func CopyFile(src, dst string) (err error) {
	srcAbs, err := filepath.Abs(src)
	if err != nil {
		return err
	}

	dstAbs, err := filepath.Abs(dst)
	if err != nil {
		return err
	}

	// open source file
	sfi, err := os.Stat(srcAbs)
	if err != nil {
		return err
	}

	if !sfi.Mode().IsRegular() {
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}

	// open dest file
	dfi, err := os.Stat(dstAbs)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}

		// file doesn't exist
		err := os.MkdirAll(filepath.Dir(dst), 0755)
		if err != nil {
			return err
		}

	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}

	return copyFileContents(src, dst)
}

//copyFileContentes copies the contents of the file named src to the file named dst
//The destination file will be created if it does not alreay exist. If the destination
//file exists, all it's contents will be replaced by the contents of the source file
func copyFileContents(src, dst string) (err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return
	}

	defer func() {
		cerr := dstFile.Close()
		if err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return
	}

	err = dstFile.Sync()
	return
}
