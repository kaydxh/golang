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

const (
	// DefaultTruncateThreshold 默认截断阈值，超过此长度才截断
	DefaultTruncateThreshold = 1024
	// DefaultTruncatePrefix 默认保留的前缀字节数
	DefaultTruncatePrefix = 10
	// MaxTruncateDepth 最大递归深度，防止循环引用（如 google.protobuf.Struct）导致堆栈溢出
	MaxTruncateDepth = 32
)

func TruncateBytes(v interface{}) interface{} {
	return TruncateBytesWithThreshold(v, DefaultTruncateThreshold, DefaultTruncatePrefix)
}

// TruncateBytesWithThreshold 截断 []byte 类型字段，超过 threshold 长度时保留前 prefix 字节并附带总长度
func TruncateBytesWithThreshold(v interface{}, threshold, prefix int) interface{} {
	return Truncate(v, func(v interface{}) bool {
		_, ok := v.([]byte)
		return ok
	}, threshold, prefix)
}

// TruncateBytesAndStrings 同时截断 []byte 和 string 类型字段
func TruncateBytesAndStrings(v interface{}) interface{} {
	return TruncateBytesAndStringsWithThreshold(v, DefaultTruncateThreshold, DefaultTruncatePrefix)
}

// TruncateBytesAndStringsWithThreshold 截断 []byte 和 string 类型字段，超过 threshold 长度时保留前 prefix 字节并附带总长度
func TruncateBytesAndStringsWithThreshold(v interface{}, threshold, prefix int) interface{} {
	return Truncate(v, func(v interface{}) bool {
		switch v.(type) {
		case []byte:
			return true
		case string:
			return true
		}
		return false
	}, threshold, prefix)
}

func Truncate(v interface{}, f func(v interface{}) bool, threshold, prefix int) interface{} {
	truncate(reflect.ValueOf(v), f, threshold, prefix, 0)
	return v
}

//https://stackoverflow.com/questions/6395076/using-reflect-how-do-you-set-the-value-of-a-struct-field
// truncate []byte, [][]byte, string, map values, not support others, eg: [][][]byte
// struct must use pointer of struct, or not rewrite it
// depth 参数用于限制递归深度，防止循环引用（如 google.protobuf.Struct）导致堆栈溢出
func truncate(v reflect.Value, f func(v interface{}) bool, threshold, prefix int, depth int) {
	if !v.IsValid() {
		return
	}

	if v.Type() == nil {
		return
	}

	// 递归深度限制，防止循环引用（如 google.protobuf.Struct -> Value -> Struct）导致堆栈溢出
	if depth > MaxTruncateDepth {
		return
	}

	if v.CanInterface() {
		vv := v.Interface()
		if f(vv) {
			truncateToLen(v, threshold, prefix)
		}
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			truncate(v.Field(i), f, threshold, prefix, depth+1)
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			truncate(v.Index(i), f, threshold, prefix, depth+1)
		}

	case reflect.Ptr:
		if !v.IsNil() {
			truncate(reflect.Indirect(v), f, threshold, prefix, depth+1)
		}

	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			mapVal := iter.Value()
			// map 的 value 不能直接修改，需要通过 Elem 解引用指针后修改
			if mapVal.Kind() == reflect.Ptr && !mapVal.IsNil() {
				truncate(mapVal.Elem(), f, threshold, prefix, depth+1)
			} else if mapVal.Kind() == reflect.Interface && !mapVal.IsNil() {
				truncate(mapVal.Elem(), f, threshold, prefix, depth+1)
			}
		}

	case reflect.Interface:
		if !v.IsNil() {
			truncate(v.Elem(), f, threshold, prefix, depth+1)
		}

	default:

	}

	return
}

func truncateToLen(oldValue reflect.Value, threshold, prefix int) {
	if !oldValue.IsValid() {
		return
	}
	if !oldValue.CanInterface() {
		return
	}

	vv := oldValue.Interface()
	switch vv := vv.(type) {
	case []byte:
		if len(vv) > threshold {
			writeBytesLenToReflectValue(oldValue, vv, prefix)
		}
	case string:
		if len(vv) > threshold {
			writeStringLenToReflectValue(oldValue, vv, prefix)
		}
	}

	return
}

// writeBytesLenToReflectValue 将 []byte 截断为前 prefix 字节 + 总长度信息
func writeBytesLenToReflectValue(v reflect.Value, data []byte, prefix int) interface{} {
	// if v can not set, return truncate result
	if !v.CanAddr() {
		return fmt.Sprintf("%s...(bytes len: %d)", string(data[:min(prefix, len(data))]), len(data))
	}

	var buf bytes.Buffer
	if prefix > 0 && prefix < len(data) {
		buf.Write(data[:prefix])
		buf.WriteString(fmt.Sprintf("...(bytes len: %d)", len(data)))
	} else {
		buf.WriteString(fmt.Sprintf("bytes len: %d", len(data)))
	}
	v.SetBytes(buf.Bytes())
	return v
}

// writeStringLenToReflectValue 将 string 截断为前 prefix 字节 + 总长度信息
func writeStringLenToReflectValue(v reflect.Value, data string, prefix int) {
	if !v.CanSet() {
		return
	}

	if prefix > 0 && prefix < len(data) {
		v.SetString(fmt.Sprintf("%s...(string len: %d)", data[:prefix], len(data)))
	} else {
		v.SetString(fmt.Sprintf("string len: %d", len(data)))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}