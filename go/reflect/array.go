/*
 *Copyright (c) 2023, kaydxh
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
	"reflect"
)

func ArrayAllTagsVaules(req interface{}, key string) []map[string]interface{} {
	sliceValues := reflect.ValueOf(req)

	var valuesMap []map[string]interface{}
	for i := 0; i < sliceValues.Len(); i++ {
		structValue := sliceValues.Index(i)
		structType := structValue.Type()

		tagValues := make(map[string]interface{})

		if structType.Kind() == reflect.Ptr {
			structType = structType.Elem()
			structValue = structValue.Elem()
		}

		if structType.Kind() != reflect.Struct {
			return nil
		}

		for j := 0; j < structType.NumField(); j++ {
			field := structType.Field(j)
			value := structValue.Field(j)
			if !value.CanInterface() {
				continue
			}

			tag := field.Tag.Get(key)
			if tag != "" {
				tagValues[tag] = value.Interface()
			}
		}

		valuesMap = append(valuesMap, tagValues)
	}

	return valuesMap
}
