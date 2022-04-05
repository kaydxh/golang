package runtime

import (
	"reflect"
	"runtime"
)

func NameOfFunction(f interface{}) string {
	v := reflect.ValueOf(f)
	if v.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	}

	return ""
}
