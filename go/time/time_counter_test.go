package time_test

import (
	"fmt"
	"testing"
	"time"

	time_ "github.com/kaydxh/golang/go/time"
)

func TestAll(t *testing.T) {
	//var tc time_.TimeCounter
	tc := time_.New(true)
	func(module string) {
		time.Sleep(time.Second)
	}("module1")
	tc.Tick("module1")

	func(module string) {
		time.Sleep(time.Second * 4)
	}("module2")
	tc.Tick("module2")

	fmt.Println(tc.String())

}
