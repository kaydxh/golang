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
package reflect_test

import (
	"testing"

	reflect_ "github.com/kaydxh/golang/go/reflect"
)

func TestReflectValue(t *testing.T) {
	testCases := []struct {
		valueType string
		value     interface{}
		expected  string
	}{
		{
			valueType: "bool",
			value:     true,
			expected:  "",
		},
		{
			valueType: "int",
			value:     123456789,
			expected:  "",
		},
		{
			valueType: "uint",
			value:     uint(123456789),
			expected:  "",
		},
		{
			valueType: "float32",
			value:     0.123456789,
			expected:  "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.valueType, func(t *testing.T) {
			rvalue, err := reflect_.ReflectValue(testCase.valueType, testCase.value)
			if err != nil {
				t.Fatalf("failed to reflect value: %v, got : %s", testCase.value, err)

			}

			t.Logf("reflect value: %v", rvalue)

		})
	}

}
