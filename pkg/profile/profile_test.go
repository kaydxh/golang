package profile_test

import (
	"fmt"
	"os"
	"testing"

	profile_ "github.com/kaydxh/golang/pkg/profile"
)

func doCpu() {
	var (
		a uint64 = 2
		i uint64
	)
	for i = 0; i < 10000000000; i++ {
		a += 2
	}
	fmt.Printf("a: %d\n", a)
}

func TestStartWithEnv(t *testing.T) {
	os.Setenv("PROFILING", "cpu")
	defer profile_.StartWithEnv().Stop()

	doCpu()
}
