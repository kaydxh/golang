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
package strings_test

import (
	"testing"

	strconv_ "github.com/kaydxh/golang/go/strconv"
	strings_ "github.com/kaydxh/golang/go/strings"
	"github.com/stretchr/testify/assert"
)

func TestReplace(t *testing.T) {
	testCases := []struct {
		s        string
		old      string
		news     []interface{}
		n        int
		expected string
	}{
		{
			s:        "task_id in (?)",
			old:      "?",
			news:     []interface{}{"a"},
			n:        1,
			expected: `task_id in ("a")`,
		},
		{
			s:        "task_id in (?, ?, ?, ?, ?)",
			old:      "?",
			news:     []interface{}{"a", "a", "a", "a"},
			n:        5,
			expected: `task_id in ("a", "a", "a", "a", "a")`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.s, func(t *testing.T) {
			newStr := strings_.Replace(testCase.s, testCase.old, testCase.news, true, testCase.n)
			t.Logf("newStr: %v", newStr)
			assert.Equal(t, testCase.expected, newStr)

		})

	}
}

func TestReplaceAll(t *testing.T) {
	testCases := []struct {
		s        string
		old      string
		news     []interface{}
		expected string
	}{
		{
			s:        "task_id in (?)",
			old:      "?",
			news:     []interface{}{"a"},
			expected: `task_id in ("a")`,
		},
		{
			s:        "task_id in (?, ?, ?, ?, ?)",
			old:      "?",
			news:     []interface{}{"a", "a", "a", "a"},
			expected: `task_id in ("a", "a", "a", "a", "a")`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.s, func(t *testing.T) {
			newStr := strings_.ReplaceAll(testCase.s, testCase.old, testCase.news, true)
			t.Logf("newStr: %v", newStr)
			assert.Equal(t, testCase.expected, newStr)

		})

	}
}

func TestSplit(t *testing.T) {
	testCases := []struct {
		s   string
		sep string
	}{
		{
			s:   "1",
			sep: ",",
		},
		{
			s:   "a,b,c",
			sep: ",",
		},
		{
			s:   "a,b,c,,,,,",
			sep: ",",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.s, func(t *testing.T) {
			ss := strings_.SplitOmitEmpty(testCase.s, testCase.sep)
			t.Logf("ss: %v. len(ss): %v", ss, len(ss))

		})

	}

}

func TestSplitToNums(t *testing.T) {
	testCases := []struct {
		s   string
		sep string
	}{
		{
			s:   "1,2,3,4",
			sep: ",",
		},
		{
			s:   "1,2,3,4,",
			sep: ",",
		},
		{
			s:   "1,2,3,4,a",
			sep: ",",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.s, func(t *testing.T) {
			nums, err := strings_.SplitToNums(testCase.s, testCase.sep, strconv_.ToInt64)
			if err != nil {
				t.Errorf("failed to split to nums, err: %v", err)
				return
			}
			t.Logf("got nums %v", nums)

		})

	}

}

func TestGetStringOrFallback(t *testing.T) {
	type args struct {
		value        string
		defaultValue string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				value:        "abc",
				defaultValue: "default",
			},
			want: "abc",
		},
		{
			name: "test2",
			args: args{
				value:        "",
				defaultValue: "default",
			},
			want: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strings_.GetStringOrFallback(tt.args.value, tt.args.defaultValue); got != tt.want {
				t.Errorf("GetStringOrFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}
