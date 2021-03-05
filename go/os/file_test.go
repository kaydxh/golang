package os_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	os_ "github.com/kaydxh/golang/go/os"
)

func TestOpenAll(t *testing.T) {
	srcDir, err := ioutil.TempDir("", "srcDir")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	defer os.RemoveAll(dstDir)

	f, err := os_.OpenAll(srcDir, os.O_CREATE, 0755)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	_ = f

	fmt.Println(srcDir)

	for i := 0; i < 10; i++ {
		dirName := filepath.Join(srcDir, fmt.Sprintf("srcdir-%d", i))
		fmt.Println(dirName)
		_, err = os_.OpenAll(dirName, os.O_CREATE, 0755)
		if err != nil {
			t.Errorf("expect nil, got %v", err)
		}
	}
}
