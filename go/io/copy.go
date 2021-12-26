/*
MIT License

Copyright (c) 2020 kay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package io

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	os_ "github.com/kaydxh/golang/go/os"
)

// Mode indicates whether to use hardlink or copy content
type Mode int

const (
	// Content creates a new file, and copies the content of the file
	Content Mode = iota
	// Hardlink creates a new hardlink to the existing file
	Hardlink
)

func CopyAll(src, dst string) (err error) {
	isDir, err := os_.IsDir(src)
	if err != nil {
		return err
	}

	if isDir {
		return CopyDir(src, dst, Content)
	}

	return CopyFile(src, dst)
}

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

func CopyRegular(srcPath, dstPath string, fileInfo os.FileInfo) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// If the destination file already exists, we shouldn't blow it away
	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, fileInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if err = doCopyWithFileClone(srcFile, dstFile); err == nil {
		return nil
	}

	if err = doCopyWithFileRange(srcFile, dstFile, fileInfo); err == nil {
		return nil
	}

	return legacyCopy(srcFile, dstFile)
}

func legacyCopy(srcFile io.Reader, dstFile io.Writer) error {
	_, err := io.Copy(dstFile, srcFile)

	return err
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
		return fmt.Errorf(
			"CopyFile: non-regular source file %s (%q)",
			sfi.Name(),
			sfi.Mode().String(),
		)
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

		if dfi.Mode().IsRegular() {
			if os.SameFile(sfi, dfi) {
				return
			}

		} else if dfi.Mode().IsDir() {
			// dst path is exist and dst is dir, so copy file to dst dir
			// src: src/1.jpg  dst: dst/dir/ => dst/dir/1.jpg
			dst = filepath.Join(dst, filepath.Base(src))

		} else {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
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
