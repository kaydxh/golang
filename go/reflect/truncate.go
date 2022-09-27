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
package reflect

import (
	"bytes"
	"fmt"
	"reflect"
)

func TruncateBytes(v interface{}) interface{} {
	return Truncate(v, func(v interface{}) bool {
		_, ok := v.([]byte)
		return ok
	})

}

func Truncate(v interface{}, f func(v interface{}) bool) interface{} {
	truncate(reflect.ValueOf(v), f)
	return v

}

//https://stackoverflow.com/questions/6395076/using-reflect-how-do-you-set-the-value-of-a-struct-field
// truncate []byte, [][]byte, not support others, eg: [][][]byte
// struct must use pointer of sturct, or not rewrite it
func truncate(v reflect.Value, f func(v interface{}) bool) {
	if !v.IsValid() {
		return
	}

	if v.Type() == nil {
		return
	}

	if v.CanInterface() {
		vv := v.Interface()
		if f(vv) {
			truncateToLen(v)
		}
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			truncate(v.Field(i), f)
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			truncate(v.Index(i), f)
		}

	case reflect.Ptr:
		truncate(reflect.Indirect(v), f)

	default:

	}

	return

}

func truncateToLen(oldValue reflect.Value) {
	if !oldValue.IsValid() {
		return
	}
	if !oldValue.CanInterface() {
		return
	}

	vv := oldValue.Interface()
	switch vv := vv.(type) {
	case []byte:
		writeLenToReflectValue(oldValue, len(vv))
	}

	return
}

func writeLenToReflectValue(v reflect.Value, length int) interface{} {
	// if v can not set, return truncate result
	if !v.CanAddr() {
		return fmt.Sprintf("bytes len: %v", length)
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("bytes len: %v", length))
	v.SetBytes(buf.Bytes())
	return v
}
