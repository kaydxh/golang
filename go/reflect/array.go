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
	values := reflect.ValueOf(req)

	var valuesMap []map[string]interface{}
	for i := 0; i < values.Len(); i++ {
		value := values.Index(i)
		if value.Kind() != reflect.Struct {
			continue
		}

		tagValues := map[string]interface{}{}
		for j := 0; j < value.NumField(); j++ {
			f := value.Field(j)
			if !f.CanInterface() {
				continue
			}

			field := value.Type().Field(j)
			if tag := field.Tag.Get(key); tag != "" {
				tagValues[tag] = f.Interface()
			}
		}

		valuesMap = append(valuesMap, tagValues)
	}

	return valuesMap
}