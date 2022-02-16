package os_test

import (
	"testing"
	"time"

	os_ "github.com/kaydxh/golang/go/os"
)

func TestGetPidsByName(t *testing.T) {
	timeout := 3000
	pids, msg, err := os_.GetPidsByName(time.Duration(timeout), "kayxhding")
	if err != nil {
		t.Errorf("expect nil, got %v, msg: %v", err, msg)
	}

	t.Logf("pids: %v", pids)
}
