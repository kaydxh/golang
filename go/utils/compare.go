package utils

import "reflect"

func Equal(a, b interface{}) bool {
	if a == nil || b == nil {
		return a == b
	}
	v1 := reflect.ValueOf(a)
	v2 := reflect.ValueOf(b)
	if v1.Type() != v2.Type() {
		return false
	}

	return a == b
}

func DeepEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}
