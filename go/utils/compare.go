package utils

import (
	"fmt"
	"reflect"
)

type Comparable interface {
	Equals(a, b interface{}) bool
}

type Comparator interface {
	Compare(a, b interface{}) int
}

func CheckComparableTypes(objs ...interface{}) bool {
	for _, obj := range objs {
		if !CheckComparableType(obj) {
			return false
		}
	}

	return true
}

func CheckComparableType(a interface{}) bool {
	switch reflect.ValueOf(a).Kind() {
	case reflect.Array, reflect.Slice, reflect.Struct, reflect.Map:
		return false

	default:
		return true
	}
}

// object (a or b) that contains members that has equality defined will
// have equality defined or die.
func Equal(a, b interface{}) (bool, error) {
	if !CheckComparableTypes(a, b) {
		_, ok := a.(Comparable)
		if !ok {
			return false, fmt.Errorf("not support compare type: %v", reflect.ValueOf(a).Kind())

		}
		_, ok = b.(Comparable)
		if !ok {
			return false, fmt.Errorf("not support compare type: %v", reflect.ValueOf(a).Kind())

		}
	}
	if a == nil || b == nil {
		return a == b, nil
	}
	v1 := reflect.ValueOf(a)
	v2 := reflect.ValueOf(b)
	if v1.Type() != v2.Type() {
		return false, nil
	}

	return a == b, nil
}

func DeepEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func Compare(a, b interface{}, compare Comparator) int {

	return compare.Compare(a, b)
}
