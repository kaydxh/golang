package syscall_test

import (
	"testing"

	syscall_ "github.com/kaydxh/golang/go/syscall"
	"gotest.tools/assert"
)

func TestGetNumFiles(t *testing.T) {

	curOpenFiles, err := syscall_.GetNumFiles()
	if err != nil {
		t.Fatalf("failed to get num files, err: %v", err)
	}

	t.Logf("curOpenFiles: %v", curOpenFiles)
}

func TestSetNumFiles(t *testing.T) {

	var maxOpenFiles uint64 = 100000
	err := syscall_.SetNumFiles(maxOpenFiles)
	if err != nil {
		t.Fatalf("failed to set num files, err: %v", err)
	}

	curOpenFiles, err := syscall_.GetNumFiles()
	if err != nil {
		t.Fatalf("failed to get num files, err: %v", err)
	}
	assert.Equal(t, maxOpenFiles, curOpenFiles)

	t.Logf("curOpenFiles: %v", curOpenFiles)
}
