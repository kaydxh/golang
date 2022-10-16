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
package strconv_test

import (
	"fmt"
	"testing"

	strconv_ "github.com/kaydxh/golang/go/strconv"
)

func TestParseNumOrDefault(t *testing.T) {

	testCases := []struct {
		str      string
		expected string
	}{
		{
			str:      "12345",
			expected: "",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			num := strconv_.ParseNumOrDefault(testCase.str, 0, strconv_.ToInt)
			t.Logf("get num: %v", num)

		})
	}

}

func TestParseNum(t *testing.T) {

	testCases := []struct {
		str      string
		expected string
	}{
		{
			str:      "12345",
			expected: "",
		},
		{
			str:      "badcase",
			expected: "",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			num, err := strconv_.ParseNum(testCase.str, strconv_.ToInt)
			if err != nil {
				t.Errorf("expecet nil, got %v", err)
			}
			t.Logf("get num: %v", num)
		})
	}

}

func TestParseNums(t *testing.T) {

	testCases := []struct {
		strs     []string
		expected string
	}{
		{
			strs:     []string{"12345"},
			expected: "",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			nums, err := strconv_.ParseNums(testCase.strs, strconv_.ToInt)
			if err != nil {
				t.Errorf("expecet nil, got %v", err)
			}
			t.Logf("get nums: %v", nums)
		})
	}

}
