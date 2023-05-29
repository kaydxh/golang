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
package math_test

import (
	"fmt"
	"testing"

	math_ "github.com/kaydxh/golang/go/math"
)

func TestEqual(t *testing.T) {
	testCases := []struct {
		a        float32
		b        float32
		expected string
	}{
		{
			a:        10.1,
			b:        10.1009,
			expected: "",
		},
		{
			a:        10.1,
			b:        10.100000009,
			expected: "",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case-%v", i), func(t *testing.T) {
			equal := math_.Equal(testCase.a, testCase.b)
			t.Logf("equal: %v", equal)
		})
	}
}
