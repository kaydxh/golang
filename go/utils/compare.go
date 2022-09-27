/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
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
