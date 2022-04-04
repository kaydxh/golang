package runtime

import (
	"reflect"
	"runtime"
)

func NameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
