package io_test

import (
	"github.com/kaydxh/golang/go/io"
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

	buf := []byte("hello world")
	err = ioutil.WriteFile(srcFilename, buf, 0777)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	err = io.CopyFile(srcFilename, dstFilename)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
}
