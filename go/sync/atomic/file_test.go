package atomic_test

import (
	"testing"

	"github.com/kaydxh/golang/go/sync/atomic"
	atomic_ "github.com/kaydxh/golang/go/sync/atomic"
)

func TestTryLock(t *testing.T) {
	var mu atomic_.File = atomic.File("test_lockfile")
	err := mu.TryLock()
	if err != nil {
		t.Errorf("expect nil, got %v", err)
		return
	}

	err = mu.TryUnLock()
	if err != nil {
		t.Errorf("expect nil, got %v", err)
		return
	}

	return
}
