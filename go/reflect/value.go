package reflect

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
)

var (
	typesMap = map[string]reflect.Type{
		// base types
		"bool":    reflect.TypeOf(true),
		"int":     reflect.TypeOf(int(1)),
		"int8":    reflect.TypeOf(int8(1)),
		"int16":   reflect.TypeOf(int16(1)),
		"int32":   reflect.TypeOf(int32(1)),
		"int64":   reflect.TypeOf(int64(1)),
		"uint":    reflect.TypeOf(uint(1)),
		"uint8":   reflect.TypeOf(uint8(1)),
		"uint16":  reflect.TypeOf(uint16(1)),
		"uint32":  reflect.TypeOf(uint32(1)),
		"uint64":  reflect.TypeOf(uint64(1)),
		"float32": reflect.TypeOf(float32(0.5)),
		"float64": reflect.TypeOf(float64(0.5)),
		"string":  reflect.TypeOf(string("")),
		// slices
		"[]bool":    reflect.TypeOf(make([]bool, 0)),
		"[]int":     reflect.TypeOf(make([]int, 0)),
		"[]int8":    reflect.TypeOf(make([]int8, 0)),
		"[]int16":   reflect.TypeOf(make([]int16, 0)),
		"[]int32":   reflect.TypeOf(make([]int32, 0)),
		"[]int64":   reflect.TypeOf(make([]int64, 0)),
		"[]uint":    reflect.TypeOf(make([]uint, 0)),
		"[]uint8":   reflect.TypeOf(make([]uint8, 0)),
		"[]uint16":  reflect.TypeOf(make([]uint16, 0)),
		"[]uint32":  reflect.TypeOf(make([]uint32, 0)),
		"[]uint64":  reflect.TypeOf(make([]uint64, 0)),
		"[]float32": reflect.TypeOf(make([]float32, 0)),
		"[]float64": reflect.TypeOf(make([]float64, 0)),
		"[]byte":    reflect.TypeOf(make([]byte, 0)),
		"[]string":  reflect.TypeOf([]string{""}),
	}

	ctxType = reflect.TypeOf((*context.Context)(nil)).Elem()

	typeConversionError = func(argValue interface{}, argTypeStr string) error {
		return fmt.Errorf("%v is not %v", argValue, argTypeStr)
	}
)

// cmd/compile/internal/gc/dump.go
func IsZeroValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return math.Float64bits(v.Float()) == 0
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		return math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !v.Index(i).IsZero() {
				return false
			}
		}
		return true
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !v.Field(i).IsZero() {
				return false
			}
		}
		return true
	default:
	}

	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

// ReflectValue converts interface{} to reflect.Value based on string type
func ReflectValue(valueType string, value interface{}) (reflect.Value, error) {
	/*
		if strings.HasPrefix(valueType, "[]") {
			return reflectValues(valueType, value)
		}
	*/

	return reflectValue(valueType, value)
}

// reflectValue converts interface{} to reflect.Value based on string type
// representing a base type (not a slice)
func reflectValue(valueType string, value interface{}) (reflect.Value, error) {
	theType, ok := typesMap[valueType]
	if !ok {
		return reflect.Value{}, NewErrUnsupportedType(valueType)
	}
	theValue := reflect.New(theType)

	// Booleans
	if theType.String() == "bool" {
		boolValue, err := getBoolValue(theType.String(), value)
		if err != nil {
			return reflect.Value{}, err
		}

		theValue.Elem().SetBool(boolValue)
		return theValue.Elem(), nil
	}

	// Integers
	if strings.HasPrefix(theType.String(), "int") {
		intValue, err := getIntValue(theType.String(), value)
		if err != nil {
			return reflect.Value{}, err
		}

		theValue.Elem().SetInt(intValue)
		return theValue.Elem(), err
	}

	// Unsigned integers
	if strings.HasPrefix(theType.String(), "uint") {
		uintValue, err := getUintValue(theType.String(), value)
		if err != nil {
			return reflect.Value{}, err
		}

		theValue.Elem().SetUint(uintValue)
		return theValue.Elem(), err
	}

	// Floating point numbers
	if strings.HasPrefix(theType.String(), "float") {
		floatValue, err := getFloatValue(theType.String(), value)
		if err != nil {
			return reflect.Value{}, err
		}

		theValue.Elem().SetFloat(floatValue)
		return theValue.Elem(), err
	}

	// Strings
	if theType.String() == "string" {
		stringValue, err := getStringValue(theType.String(), value)
		if err != nil {
			return reflect.Value{}, err
		}

		theValue.Elem().SetString(stringValue)
		return theValue.Elem(), nil
	}

	return reflect.Value{}, NewErrUnsupportedType(valueType)
}

func getBoolValue(theType string, value interface{}) (bool, error) {
	b, ok := value.(bool)
	if !ok {
		return false, typeConversionError(value, typesMap[theType].String())
	}

	return b, nil
}

func getIntValue(theType string, value interface{}) (int64, error) {
	// We use https://golang.org/pkg/encoding/json/#Decoder.UseNumber when unmarshaling signatures.
	// This is because JSON only supports 64-bit floating point numbers and we could lose precision
	// when converting from float64 to signed integer
	if strings.HasPrefix(fmt.Sprintf("%T", value), "json.Number") {
		n, ok := value.(json.Number)
		if !ok {
			return 0, typeConversionError(value, typesMap[theType].String())
		}

		return n.Int64()
	}

	var n int64
	switch value := value.(type) {
	case int:
		n = int64(value)
	case int64:
		n = value
	case int32:
		n = int64(value)
	case int16:
		n = int64(value)
	case int8:
		n = int64(value)
	default:
		fmt.Printf("value: %v\n", value)
		return 0, typeConversionError(value, typesMap[theType].String())
	}

	return n, nil
}

func getUintValue(theType string, value interface{}) (uint64, error) {
	// We use https://golang.org/pkg/encoding/json/#Decoder.UseNumber when unmarshaling signatures.
	// This is because JSON only supports 64-bit floating point numbers and we could lose precision
	// when converting from float64 to unsigned integer
	if strings.HasPrefix(fmt.Sprintf("%T", value), "json.Number") {
		n, ok := value.(json.Number)
		if !ok {
			fmt.Printf("00000\n")
			return 0, typeConversionError(value, typesMap[theType].String())
		}

		intVal, err := n.Int64()
		if err != nil {
			return 0, err
		}

		return uint64(intVal), nil
	}

	var n uint64
	switch value := value.(type) {
	case uint:
		n = uint64(value)
	case uint64:
		n = value
	case uint32:
		n = uint64(value)
	case uint16:
		n = uint64(value)
	case uint8:
		n = uint64(value)
	default:
		fmt.Printf("value: %v\n", value)
		return 0, typeConversionError(value, typesMap[theType].String())
	}
	return n, nil
}

func getFloatValue(theType string, value interface{}) (float64, error) {
	// We use https://golang.org/pkg/encoding/json/#Decoder.UseNumber when unmarshaling signatures.
	// This is because JSON only supports 64-bit floating point numbers and we could lose precision
	if strings.HasPrefix(fmt.Sprintf("%T", value), "json.Number") {
		n, ok := value.(json.Number)
		if !ok {
			return 0, typeConversionError(value, typesMap[theType].String())
		}

		return n.Float64()
	}

	f, ok := value.(float64)
	if !ok {
		return 0, typeConversionError(value, typesMap[theType].String())
	}

	return f, nil
}

func getStringValue(theType string, value interface{}) (string, error) {
	s, ok := value.(string)
	if !ok {
		return "", typeConversionError(value, typesMap[theType].String())
	}

	return s, nil
}
