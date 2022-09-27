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
package url

import (
	"reflect"
	"strconv"
)

var (
	boolV      = boolValue{}
	intV       = intValue{}
	unitV      = uintValue{}
	floatV     = floatValue{}
	stringV    = stringValue{}
	encoderMap = map[reflect.Kind]Value{
		reflect.Bool:    boolV,
		reflect.Int:     intV,
		reflect.Int8:    intV,
		reflect.Int16:   intV,
		reflect.Int32:   intV,
		reflect.Int64:   intV,
		reflect.Uint:    unitV,
		reflect.Uint8:   unitV,
		reflect.Uint16:  unitV,
		reflect.Uint32:  unitV,
		reflect.Uint64:  unitV,
		reflect.Uintptr: unitV,
		reflect.Float32: floatV,
		reflect.Float64: floatV,
		reflect.String:  stringV,
	}
)

type Value interface {
	Encode(value reflect.Value) string
	//Decode(value string) (reflect.Value, error)
}

type boolValue struct{}

func (e boolValue) Encode(value reflect.Value) string {
	if value.Bool() {
		return "1"
	}
	return "0"
}

type intValue struct{}

func (e intValue) Encode(value reflect.Value) string {
	return strconv.FormatInt(value.Int(), 10)
}

type uintValue struct{}

func (e uintValue) Encode(value reflect.Value) string {
	return strconv.FormatUint(value.Uint(), 10)
}

type floatValue struct{}

func (e floatValue) Encode(value reflect.Value) string {
	return strconv.FormatFloat(value.Float(), 'f', -1, 64)
}

type stringValue struct{}

func (e stringValue) Encode(value reflect.Value) string {
	return value.String()
}

func getEncoder(kind reflect.Kind) Value {
	if encoder, ok := encoderMap[kind]; ok {
		return encoder
	}

	return nil
}
