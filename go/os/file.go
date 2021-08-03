package os

import (
	"fmt"
	"os"
	"path/filepath"
)

func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err

}

// MkdirAll creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
// The permission(0755) bits perm (before umask) are used for all
// directories that MkdirAll creates.
// If path is already a directory, MkdirAll does nothing
func MakeDirAll(name string) error {
	return os.MkdirAll(name, 0755)
}

func MakeDir(name string) error {
	return os.Mkdir(name, 0755)
}

func OpenAll(path string, flag int, perm os.FileMode) (*os.File, error) {
	dir, file := filepath.Split(path)
	// file or dir exists
	if _, err := os.Stat(path); err == nil {
		return os.OpenFile(path, flag, perm)
	}

	// mkdir -p dir
	if err := os.MkdirAll(dir, perm); err != nil {
		return nil, err
	}

	if file == "" {
		return nil, nil
	}

	return os.OpenFile(path, flag, perm)
}

func OpenFile(path string, appended bool) (file *os.File, err error) {
	if !appended {
		file, err = OpenAll(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	} else {
		file, err = OpenAll(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	}

	return file, err

}

// SameFile reports whether fi1 and fi2 describe the same file.
func SameFile(fi1, fi2 string) bool {
	stat1, err := os.Stat(fi1)
	if err != nil {
		return false
	}

	stat2, err := os.Stat(fi2)
	if err != nil {
		return false
	}
	return os.SameFile(stat1, stat2)
}

//oldname and newname is full path
func SymLink(oldname, newname string) error {

	_, err := os.Stat(oldname)
	if err != nil {
		return fmt.Errorf("failed to stat oldname: %v, err: %v", oldname, err)
	}

	linkdir := filepath.Dir(newname)
	err = MakeDirAll(linkdir)
	if err != nil {
		return fmt.Errorf("failed to make dir: %v, err: %v", linkdir, err)
	}

	// check link file is valid
	_, err = os.Readlink(newname)
	if err == nil {
		return nil
	}

	//link file is invalid
	os.Remove(newname)
	/*
		_, err = os.Stat(newname)
		if err != nil {
			if os.IsNotExist(err) {
				os.Remove(newname)
			}
		}
	*/
	err = os.Symlink(oldname, newname)
	if err != nil {
		return fmt.Errorf("failed to symlink err: %v", err)
	}

	_, err = os.Stat(newname)
	if err != nil {
		return fmt.Errorf("failed to stat newname: %v, err: %v", newname, err)
	}

	return err
}
