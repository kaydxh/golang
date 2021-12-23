package reflect

import (
	"bytes"
	"fmt"
	"reflect"
)

//https://stackoverflow.com/questions/6395076/using-reflect-how-do-you-set-the-value-of-a-struct-field
// truncate []byte, [][]byte, not support others, eg: [][][]byte
// struct must use pointer of sturct, or not rewrite it
func TruncateBytes(req interface{}) interface{} {
	v, ok := indirectStruct(req)
	if !ok {
		return nil
	}
	if !v.IsValid() {
		return nil
	}

	tt := reflect.TypeOf(req)
	if tt.Kind() == reflect.Ptr {
		tt = tt.Elem()
	}
	if tt.Kind() != reflect.Struct {
		return truncateToLen(v)
	}

	for i := 0; i < tt.NumField(); i++ {
		field := tt.Field(i)
		property := string(field.Name)
		f := v.FieldByName(property)
		if !f.IsValid() {
			return nil
		}

		if !f.CanSet() {
			continue
		}
		if !f.CanInterface() {
			continue
		}

		valueValue := v.Field(i).Addr()
		switch v.Field(i).Kind() {
		case reflect.Struct:
			TruncateBytes(valueValue.Interface())
		}
		truncateToLen(f)
	}

	return req
}

func truncateToLen(oldValue reflect.Value) interface{} {
	if !oldValue.IsValid() {
		return oldValue
	}
	if !oldValue.CanInterface() {
		return oldValue
	}

	vv := oldValue.Interface()
	switch vv := vv.(type) {
	case [][]byte:
		for i, subV := range vv {
			writeLenToReflectValue(oldValue.Index(i), len(subV))
		}

	case []byte:
		return writeLenToReflectValue(oldValue, len(vv))
	}

	return oldValue
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
