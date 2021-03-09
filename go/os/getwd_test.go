package os_test

import (
	"fmt"
	"testing"

	os_ "github.com/kaydxh/golang/go/os"
)

func TestGetwd(t *testing.T) {
	path, err := os_.Getwd()
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	fmt.Println(path)
}
