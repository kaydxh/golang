package os_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

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
