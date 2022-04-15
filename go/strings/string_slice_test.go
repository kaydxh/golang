package strings_test

import (
	"fmt"
	"testing"

	strings_ "github.com/kaydxh/golang/go/strings"
	"github.com/stretchr/testify/assert"
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

func TestRemoveEmpty(t *testing.T) {
	testCases := []struct {
		s []string
	}{
		{
			s: []string{"1", "", "3", "4"},
		},
		{
			s: []string{"", "", "", ""},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			intersection := strings_.RemoveEmpty(testCase.s)
			t.Logf("intersection :%v", intersection)
		})
	}
}

func TestSliceContains(t *testing.T) {
	testCases := []struct {
		s             []string
		target        string
		caseSensitive bool
		expected      bool
	}{
		{
			s:             []string{"a", "bdx"},
			target:        "bDx",
			caseSensitive: false,
			expected:      true,
		},
		{
			s:             []string{"a", "bdx"},
			target:        "bdx",
			caseSensitive: true,
			expected:      true,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			has := strings_.SliceContains(testCase.s, testCase.target, testCase.caseSensitive)
			t.Logf("has :%v", has)
			assert.Equal(t, testCase.expected, has)
		})
	}
}
