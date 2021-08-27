package strings_test

import (
	"testing"

	strings_ "github.com/kaydxh/golang/go/strings"
)

func TestStringIntersection(t *testing.T) {
	testCases := []struct {
		name     string
		s1       []string
		s2       []string
		expected []string
	}{
		{
			name:     "test string",
			s1:       []string{"1", "2", "3", "4"},
			s2:       []string{"3", "4", "5", "6"},
			expected: []string{"3", "4"},
		},
		{
			name:     "test string2",
			s1:       []string{"1", "2", "3", "4"},
			s2:       []string{"5", "6"},
			expected: []string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			intersection := strings_.SliceIntersection(testCase.s1, testCase.s2)
			if len(intersection) != len(testCase.expected) {
				t.Fatalf("Expected Intersection len: %v, got : %v", len(testCase.expected), len(intersection))
			}
			t.Logf("intersection :%v", intersection)
		})
	}
}
