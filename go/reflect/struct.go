package reflect

import "reflect"

func indirectStruct(req interface{}) (reflect.Value, bool) {
	if req == nil {
		return reflect.Value{}, false
	}
	v := reflect.ValueOf(req)
	if v.IsNil() {
		return reflect.Value{}, false
	}

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
