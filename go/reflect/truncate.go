package reflect

import (
	"bytes"
	"fmt"
	"reflect"
)

//https://stackoverflow.com/questions/6395076/using-reflect-how-do-you-set-the-value-of-a-struct-field
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
		vv, ok := v.Interface().([]byte)
		if ok {
			req = fmt.Sprintf("bytes len: %v", len(vv))
		}
		return req
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

func truncateToLen(oldValue reflect.Value) {
	if !oldValue.IsValid() {
		return
	}
	if !oldValue.CanSet() {
		return
	}
	if !oldValue.CanInterface() {
		return
	}

	vv := oldValue.Interface()
	switch vv := vv.(type) {
	case []byte:
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("bytes len: %v", len(vv)))
		oldValue.SetBytes(buf.Bytes())
	}
}
