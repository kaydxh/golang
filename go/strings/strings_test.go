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