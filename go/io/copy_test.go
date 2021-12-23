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

package io_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/kaydxh/golang/go/io"
	io_ "github.com/kaydxh/golang/go/io"
	"gotest.tools/v3/assert"
)

func TestCopyDir(t *testing.T) {
	srcDir, err := ioutil.TempDir("", "srcDir")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	populateSrcDir(t, srcDir, 3)

	dstDir, err := ioutil.TempDir("", "dstDir")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	defer os.RemoveAll(dstDir)

	err = io.CopyDir(srcDir, dstDir, io.Content)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	err = filepath.Walk(srcDir, func(srcPath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Rebase path
		relPath, err := filepath.Rel(srcDir, srcPath)
		assert.NilError(t, err)
		if relPath == "." {
			return nil
		}

		dstPath := filepath.Join(dstDir, relPath)
		// If we add non-regular dirs and files to the test
		// then we need to add more checks here.
		dstFileInfo, err := os.Lstat(dstPath)
		assert.NilError(t, err)

		fmt.Println("dstFileInfo: ", dstFileInfo)

		return nil
	})

	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

}

func TestCopyFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "file-copy-check")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	defer os.RemoveAll(dir)

	srcFilename := filepath.Join(dir, "srcFilename")
	dstFilename := filepath.Join(dir, "dstilename")
	fmt.Println("srcFilename: , dstFilename: ", srcFilename, dstFilename)

	buf := []byte("hello world")
	err = ioutil.WriteFile(srcFilename, buf, 0777)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	err = io.CopyFile(srcFilename, dstFilename)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	readBuf, err := ioutil.ReadFile(srcFilename)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	fmt.Println("srcFilename content: ", string(readBuf))

	readBuf, err = ioutil.ReadFile(dstFilename)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	fmt.Println("dstFilename content: ", string(readBuf))

	if !bytes.Equal(buf, readBuf) {
		t.Errorf("expect true, got %v", false)
	}
}

func TestCopyAll(t *testing.T) {
	testCases := []struct {
		src string
		dst string
	}{
		{
			src: "./testdata/file/1.txt",
			dst: "./testdata.copy/file/1.txt",
		},
		/*
			{
				src: "./testdata/dir",
				dst: "./testdata.copy/dir",
			},
		*/
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			err := io_.CopyAll(testCase.src, testCase.dst)
			if err != nil {
				t.Fatalf("copy all expected nil, got %v", err)
			}
		})
	}

}

func randomMode(baseMode int) os.FileMode {
	for i := 0; i < 7; i++ {
		baseMode = baseMode | (1&rand.Intn(2))<<uint(i)
	}
	return os.FileMode(baseMode)
}

func populateSrcDir(t *testing.T, srcDir string, remainingDepth int) {
	if remainingDepth == 0 {
		return
	}

	for i := 0; i < 10; i++ {
		dirName := filepath.Join(srcDir, fmt.Sprintf("srcdir-%d", i))
		// Owner all bits set
		assert.NilError(t, os.Mkdir(dirName, randomMode(0700)))
		populateSrcDir(t, dirName, remainingDepth-1)
	}

	for i := 0; i < 10; i++ {
		fileName := filepath.Join(srcDir, fmt.Sprintf("srcfile-%d", i))
		// Owner read bit set
		assert.NilError(t, ioutil.WriteFile(fileName, []byte{}, randomMode(0400)))
	}

}
