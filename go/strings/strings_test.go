package strings_test

import (
	"testing"

	strings_ "github.com/kaydxh/golang/go/strings"
	"github.com/stretchr/testify/assert"
)

func TestWriteFile(t *testing.T) {
	testCases := []struct {
		s        string
		old      string
		news     []string
		n        int
		expected string
	}{
		{
			s:        "task_id in (?)",
			old:      "?",
			news:     []string{`"a"`},
			n:        1,
			expected: `task_id in ("a")`,
		},
		{
			s:        "task_id in (?, ?, ?, ?, ?)",
			old:      "?",
			news:     []string{`"a"`, `"a"`, `"a"`, `"a"`},
			n:        5,
			expected: `task_id in ("a", "a", "a", "a", "a")`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.s, func(t *testing.T) {
			newStr := strings_.Replace(testCase.s, testCase.old, testCase.news, testCase.n)
			t.Logf("newStr: %v", newStr)
			assert.Equal(t, testCase.expected, newStr)

		})

	}
}
