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

	if !v.IsValid() {
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

	if !v.IsValid() {
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
	var tags []string

	tagsValues := fieldTagsValues(req, key, nonzero)
	for tag := range tagsValues {
		tags = append(tags, tag)
	}

	return tags
}

// only get export Fields
func AllTagsValues(req interface{}, key string) map[string]interface{} {
	return fieldTagsValues(req, key, false)
}

// req must be struct(Not pointer to struct), or return nil(tt.Field() will panic)
// key for tag , db or json ..., if key is empty or tag is empty, ignore it
// nonzere true, noly return field tags for values that nonzero
func fieldTagsValues(req interface{}, key string, nonzero bool) map[string]interface{} {
	if req == nil {
		return nil
	}

	v, ok := indirectStruct(req)
	if !ok {
		return nil
	}
	if !v.IsValid() {
		return nil
	}

	tt := reflect.TypeOf(req)
	if tt.Kind() != reflect.Struct {
		return nil
	}

	tagsValues := make(map[string]interface{})
	for i := 0; i < tt.NumField(); i++ {
		field := tt.Field(i)
		property := string(field.Name)
		f := v.FieldByName(property)

		if !f.CanInterface() {
			continue
		}

		if nonzero {
			if IsZeroValue(f) {
				continue
			}
		}

		if key == "" {
			continue
		}

		tag := field.Tag.Get(key)
		if len(tag) > 0 {
			// field.Type.Name() -> "string", "int64" ...
			tagsValues[tag] = f.Interface()
		}
	}

	return tagsValues
}
