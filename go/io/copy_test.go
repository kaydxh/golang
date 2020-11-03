package io_test

import (
	"bytes"
	"fmt"
	"github.com/kaydxh/golang/go/io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

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
