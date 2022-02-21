package strings_test

import (
	"testing"

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
			s:   "",
			sep: ",",
		},
		{
			s:   "a,b,c",
			sep: ",",
		},
		{
			s:   "a,b,c,",
			sep: ",",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.s, func(t *testing.T) {
			ss := strings_.Split(testCase.s, testCase.sep)
			t.Logf("ss: %v. len(ss): %v", ss, len(ss))

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
