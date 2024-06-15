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
package bytes

import (
	"reflect"
	"unsafe"
)

func Truncate(s []byte, n int) []byte {
	if n < 0 {
		n = 0
	}

	if len(s) <= n {
		return s
	}

	return s[:n]
}

func Encode[T any](m T) []byte {
	if reflect.ValueOf(m).IsZero() {
		return nil
	}

	sz := int(unsafe.Sizeof(m))
	p := make([]byte, sz)

	slice := (*reflect.SliceHeader)(unsafe.Pointer(&p))
	slice.Data = uintptr(unsafe.Pointer(&m))
	slice.Len = sz
	slice.Cap = sz
	return p
}

func Decode[T any](b []byte) T {
	var zero T
	if len(b) == 0 {
		return zero
	}

	var m T = *(*T)(unsafe.Pointer(&b[0]))
	return m
}
