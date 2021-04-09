package os

import (
	"fmt"
	"testing"

	os_ "github.com/kaydxh/golang/go/os"
)

func TestGetPidByName(t *testing.T) {
	pids, err := os_.GetPidByName("kayxhding")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	fmt.Println(pids)
}
