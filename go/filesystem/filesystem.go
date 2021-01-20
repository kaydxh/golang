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

package filesystem

import (
	"os"
	"path/filepath"
	"time"
)

// Filesystem is an interface that we can use to mock various filesystem operations
type Filesystem interface {
	// from "os"
	Stat(name string) (os.FileInfo, error)
	Create(name string) (File, error)
	Rename(oldpath, newpath string) error
	MkdirAll(path string, perm os.FileMode) error
	Chtimes(name string, atime time.Time, mtime time.Time) error
	RemoveAll(path string) error
	Remove(name string) error

	// from "io/ioutil"
	ReadFile(filename string) ([]byte, error)
	TempDir(dir, prefix string) (string, error)
	TempFile(dir, prefix string) (File, error)
	ReadDir(dirname string) ([]os.FileInfo, error)
	Walk(root string, walkFn filepath.WalkFunc) error
}

// File is an interface that we can use to mock various filesystem operations typically
// accessed through the File object from the "os" package
type File interface {
	// for now, the only os.File methods used are those below, add more as necessary
	Name() string
	Write(b []byte) (n int, err error)
	Sync() error
	Close() error
}
