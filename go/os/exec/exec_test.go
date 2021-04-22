package exec_test

import (
	"fmt"
	"testing"
	"time"

	exec_ "github.com/kaydxh/golang/go/os/exec"
)

func TestExec(t *testing.T) {
	cmd := fmt.Sprintf(`ps -elf`)
	timeout := 5000 //ms
	result, msg, err := exec_.Exec(time.Duration(timeout), "/bin/sh", "-c", cmd)
	if err != nil {
		t.Errorf("expect nil, got %v, msg: %v", err, msg)
	}
	t.Logf("result: %v", result)
}
