package slices_test

import (
	"fmt"
	"testing"

	slices_ "github.com/kaydxh/golang/go/slices"
)

func TestSliceIntersection(t *testing.T) {
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
			intersection := slices_.SliceIntersection(testCase.s1, testCase.s2)
			if len(intersection) != len(testCase.expected) {
				t.Fatalf("Expected Intersection len: %v, got : %v", len(testCase.expected), len(intersection))
			}
			t.Logf("intersection :%v", intersection)
		})
	}
}

func TestSliceDifference(t *testing.T) {
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
			expected: []string{"1", "2"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			difference := slices_.SliceDifference(testCase.s1, testCase.s2)
			if len(difference) != len(testCase.expected) {
				t.Fatalf("Expected Difference len: %v, got : %v", len(testCase.expected), len(difference))
			}
			t.Logf("difference :%v", difference)
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
			intersection := slices_.RemoveEmpty(testCase.s)
			t.Logf("intersection :%v", intersection)
		})
	}
}
