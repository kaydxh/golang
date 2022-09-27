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
			difference := strings_.SliceDifference(testCase.s1, testCase.s2)
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
