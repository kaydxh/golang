package runtime_test

import (
	"testing"

	runtime_ "github.com/kaydxh/golang/go/runtime"
)

func TestGetCaller(t *testing.T) {
	// fn shoud be testing.tRunner, caller of TestGetParentCaller
	fn := runtime_.GetCaller()
	t.Logf("fn: %v", fn)

}
