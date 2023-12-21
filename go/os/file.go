/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
package os

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fi.IsDir(), nil
}

func IsHidden(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	base := filepath.Base(path)
	if base == "." || base == ".." {
		return false, nil
	}
	if len(base) > 0 && base[0] == '.' {
		return true, nil
	}

	return false, nil
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

// if name is empty, create a directory named name, or creates a new temporary directory in the directory dir
// pattern can be empty
func MakeTempDirAll(name, pattern string) (string, error) {
	if name != "" {
		// make base dir, or if the dir is not exist, MakeTempDirAll return error
		err := MakeDirAll(name)
		if err != nil {
			return "", err
		}
	}

	return os.MkdirTemp(name, pattern)
}

func OpenAll(path string, flag int, perm os.FileMode) (*os.File, error) {
	dir, file := filepath.Split(path)
	// file or dir exists
	if _, err := os.Stat(path); err == nil {
		return os.OpenFile(path, flag, perm)
	}

	// mkdir -p dir
	if dir != "" {
		if err := MakeDirAll(dir); err != nil {
			return nil, err
		}
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

// oldname and newname is full path
func SymLink(oldname, newname string) error {

	oldname, err := filepath.Abs(oldname)
	if err != nil {
		return err
	}
	newname, err = filepath.Abs(newname)
	if err != nil {
		return err
	}

	_, err = os.Stat(oldname)
	if err != nil {
		return fmt.Errorf("failed to stat oldname: %v, err: %v", oldname, err)
	}

	linkdir := filepath.Dir(newname)
	err = MakeDirAll(linkdir)
	if err != nil {
		return fmt.Errorf("failed to make dir: %v, err: %v", linkdir, err)
	}

	// check link file is valid
	targetname, err := os.Readlink(newname)
	if err == nil {
		if targetname == oldname {
			return nil
		}
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

	_, err = os.Lstat(newname)
	if err != nil {
		return fmt.Errorf("failed to stat newname: %v, err: %v", newname, err)
	}

	return err
}

// include subdir, file and hidden file
func ReadDirNames(path string, filterHiddenFile bool) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	names, err := file.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	if !filterHiddenFile {
		return names, nil
	}

	noHiddenFiles := []string{}

	for _, name := range names {
		if strings.HasPrefix(filepath.Base(name), ".") {
			continue
		}

		noHiddenFiles = append(noHiddenFiles, name)

	}

	return noHiddenFiles, nil

}

// only inlcude subdir
// filterHiddenFile  is true, then filter hidden files from result
func ReadDirSubDirNames(path string, filterHiddenFile bool) ([]string, error) {
	names, err := ReadDirNames(path, filterHiddenFile)
	if err != nil {
		return nil, err
	}

	dirs := []string{}
	for _, name := range names {

		filePath := filepath.Join(path, name)
		ok, _ := IsDir(filePath)
		if ok {
			dirs = append(dirs, name)
		}
	}

	return dirs, nil
}

func ReadDirFileNames(path string, filterHiddenFile bool) ([]string, error) {
	names, err := ReadDirNames(path, filterHiddenFile)
	if err != nil {
		return nil, err
	}

	dirs := []string{}
	for _, name := range names {

		filePath := filepath.Join(path, name)
		ok, err := IsDir(filePath)
		if err == nil && !ok {
			dirs = append(dirs, name)
		}
	}

	return dirs, nil
}
