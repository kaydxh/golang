package runtime_test

import (
	"testing"

	runtime_ "github.com/kaydxh/golang/go/runtime"
)

func TestRecover(t *testing.T) {
	defer runtime_.Recover()
	panic("test panic")
}

func testRecoverFromPanic() (err error) {
	defer runtime_.RecoverFromPanic(&err)
	panic("test panic")
}

func TestRecoverFromPanic(t *testing.T) {
	err := testRecoverFromPanic()
	t.Logf("err: %v", err)
}
