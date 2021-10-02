package syscall_test

import (
	"testing"

	syscall_ "github.com/kaydxh/golang/go/syscall"
)

func TestGetNumFiles(t *testing.T) {

	curOpenFiles, maxOpenFiles, err := syscall_.GetNumFiles()
	if err != nil {
		t.Fatalf("failed to get num files, err: %v", err)
	}

	t.Logf("curOpenFiles: %v,  maxOpenFiles: %v", curOpenFiles, maxOpenFiles)
}

func TestSetMaxNumFiles(t *testing.T) {

	curOpenFiles, err := syscall_.SetMaxNumFiles()
	if err != nil {
		t.Fatalf("failed to set max num files, err: %v", err)
	}

	t.Logf("curOpenFiles: %v", curOpenFiles)
}
