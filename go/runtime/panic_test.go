package runtime_test

import (
	"testing"

	runtime_ "github.com/kaydxh/golang/go/runtime"
)

func TestRecover(t *testing.T) {
	defer runtime_.Recover()
	panic("test panic")
}
