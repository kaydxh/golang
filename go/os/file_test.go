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
package os_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	io_ "github.com/kaydxh/golang/go/io"
	os_ "github.com/kaydxh/golang/go/os"
	"github.com/stretchr/testify/assert"
)

func TestOpenAll(t *testing.T) {
	dir, err := ioutil.TempDir("", "dir")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	defer os.RemoveAll(dir)

	f, err := os_.OpenAll(dir, os.O_CREATE, 0755)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	_ = f

	fmt.Println(dir)

	for i := 0; i < 10; i++ {
		dirName := filepath.Join(dir, fmt.Sprintf("dir-%d", i))
		f, err = os_.OpenAll(dirName, os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			t.Errorf("expect nil, got %v", err)
		}
		defer f.Close()

		n, err := io.WriteString(f, "hello world")
		if err != nil {
			t.Errorf("expect nil, got %v", err)
		}
		assert.Equal(t, len("hello world"), n)
		fmt.Println(dirName)
	}
}

func TestSameFile(t *testing.T) {
	total := 2
	fileNames := make([]string, total)
	dir, err := ioutil.TempDir("", "dir")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	defer os.RemoveAll(dir)

	for i := 0; i < total; i++ {
		fileNames[i] = filepath.Join(dir, fmt.Sprintf("dir-%d", i))

		buf := []byte("hello world")
		err = ioutil.WriteFile(fileNames[i], buf, 0777)
		if err != nil {
			t.Errorf("expect nil, got %v", err)
		}
		fmt.Println(fileNames[i])
	}

	ok := os_.SameFile(fileNames[0], fileNames[0])
	assert.Truef(t, ok, "fileName %s is the same of %s", fileNames[0], fileNames[0])

	ok = os_.SameFile(fileNames[0], fileNames[1])
	assert.Falsef(t, ok, "fileName %s is not the same of %s", fileNames[0], fileNames[1])

}

func TestSymLink(t *testing.T) {
	dir, err := ioutil.TempDir("", "dir")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	defer os.RemoveAll(dir)

	oldname := filepath.Join(dir, "oldname")

	buf := []byte("hello world")
	err = ioutil.WriteFile(oldname, buf, 0777)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	newname := "test.link"
	err = os_.SymLink(oldname, newname)
	if err != nil {
		t.Fatalf("expect nil, got %v", err)
	}

	data, err := io_.ReadFile(newname)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	assert.Equal(t, buf, data)

}

func TestPath(t *testing.T) {

	filename := "./workspace/studyspace/.././11.jpg"
	///Users/kayxhding/workspace/studyspace/git-kayxhding/github.com/kaydxh/golang/go/os/workspace/11.jpg
	fmt.Println(filepath.Abs(filename))
	//workspace/11.jpg
	fmt.Println(filepath.Clean(filename))
	//../11.jpg
	fmt.Println(filepath.Rel("./workspace/22", filename))
}

func TestReadDirNames(t *testing.T) {
	dir := "/usr/local"
	allFiles, err := os_.ReadDirNames(dir, true)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	t.Logf("got sub dirs: %v", allFiles)

}

func TestReadSubDirs(t *testing.T) {
	dir := "/usr/local"
	subDirs, err := os_.ReadDirSubDirNames(dir, true)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	t.Logf("got sub dirs: %v", subDirs)
}

func TestReadDirFileNames(t *testing.T) {
	dir := "/usr/local"
	fileNames, err := os_.ReadDirFileNames(dir, true)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	t.Logf("got filenames: %v", fileNames)

}

func TestMakeTempDirAll(t *testing.T) {
	//  /var/folders/by/jt3bf2n52fv6v7lrbndtdmjr0000gn/T/example690523200.12
	dir, err := os_.MakeTempDirAll("", "example*.12")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	t.Logf("got dir: %v", dir)

}
