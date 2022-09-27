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
package base64_test

import (
	"testing"

	base64_ "github.com/kaydxh/golang/go/encoding/base64"
	"gotest.tools/v3/assert"
)

func TestString(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
	}{
		{
			name:     "hello word",
			expected: "",
		},
		{
			name:     "http://12306.com?a=%b",
			expected: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			encoded := base64_.EncodeString(testCase.name)
			t.Logf("base64 encode : %v", encoded)
			decoded, err := base64_.DecodeString(encoded)
			if err != nil {
				t.Fatalf("failed to decode string, err: %v", err)
			}

			assert.Equal(t, testCase.name, decoded)
		})
	}
	/*
		content := "hello word"

		encoded := base64_.EncodeString(content)
		t.Logf("base64 encode : %v", encoded)

		decoded, err := base64_.DecodeString(encoded)
		if err != nil {
			t.Fatalf("failed to decode string, err: %v", err)
		}
		t.Logf("base64 decode : %v", decoded)
	*/
	//assert.Equal(t, content, decoded)
}

func TestURL(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
	}{
		{
			name:     `http://baidu.com?a=10&b="hello"`,
			expected: "",
		},
		{
			name:     "http://12306.com?a=%b",
			expected: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			encoded := base64_.EncodeURL(testCase.name)
			t.Logf("base64 encode : %v", encoded)
			decoded, err := base64_.DecodeURL(encoded)
			if err != nil {
				t.Fatalf("failed to decode url, err: %v", err)
			}

			assert.Equal(t, testCase.name, decoded)
		})
	}
}
