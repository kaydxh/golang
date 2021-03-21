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
