package reflect

import (
	"reflect"
)

//req must be pointer to struct
// IsNil reports whether its argument v is nil. The argument must be
// a chan, func, interface, map, pointer, or slice value; if it is
// not, IsNil panics. Note that IsNil is not always equivalent to a
// regular comparison with nil in Go. For example, if v was created
// by calling ValueOf with an uninitialized interface variable i,
// i==nil will be true but v.IsNil will panic as v will be the zero
// Value.
func indirectStruct(req interface{}) (reflect.Value, bool) {
	if req == nil {
		return reflect.Value{}, false
	}

	v := reflect.ValueOf(req)
	/*
		if v.IsNil() {
			return reflect.Value{}, false
		}
	*/

	return reflect.Indirect(v), true
}

func RetrieveStructField(req interface{}, name string) string {
	v, ok := indirectStruct(req)
	if !ok {
		return ""
	}

	//nested field: reflect.Indirect(v).FieldByName("layout1").Index(0).FieldByName("layout2")
	f := v.FieldByName(name)
	if f.IsValid() && f.Kind() == reflect.String {
		return f.String()
	}
	return ""

}

func TrySetStructFiled(req interface{}, name, value string) {
	v, ok := indirectStruct(req)
	if !ok {
		return
	}
	f := v.FieldByName(name)
	if f.IsValid() && f.Kind() == reflect.String {
		f.SetString(value)
	}
}

func AllFieldTags(req interface{}, key string) []string {
	return fieldTags(req, key, false)
}

// req must be struct(Not pointer to struct), or return nil(tt.Field() will panic)
// key for tag , db or json, if key is empty, use field name instead
func NonzeroFieldTags(req interface{}, key string) []string {
	return fieldTags(req, key, true)
}

// req must be struct(Not pointer to struct), or return nil(tt.Field() will panic)
// key for tag , db or json, if key is empty, use field name instead
//nonzere true, noly return field tags for values that nonzero
func fieldTags(req interface{}, key string, nonzero bool) []string {
	if req == nil {
		return nil
	}

	v, ok := indirectStruct(req)
	if !ok {
		return nil
	}

	tt := reflect.TypeOf(req)
	if tt.Kind() != reflect.Struct {
		return nil
	}

	var fields []string
	for i := 0; i < tt.NumField(); i++ {
		property := string(tt.Field(i).Name)
		f := v.FieldByName(property)
		if nonzero {
			if IsZeroValue(f) {
				continue
			}
		}

		if len(key) == 0 {
			fields = append(fields, property)
		} else {
			fields = append(fields, string(tt.Field(i).Tag.Get(key)))
		}

	}

	return fields
}
