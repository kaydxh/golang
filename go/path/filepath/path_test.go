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
package filepath_test

import (
	"fmt"
	"path/filepath"
	"testing"

	filepath_ "github.com/kaydxh/golang/go/path/filepath"
	"github.com/stretchr/testify/assert"
)

func TestGetParentRelPath(t *testing.T) {
	testCases := []struct {
		filePath string
		expected string
	}{
		{
			filePath: "/data/home/deploy/bin/run.exe",
			expected: "bin/run.exe",
		},
		{
			filePath: "/run.exe",
			expected: "/run.exe",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.filePath, func(t *testing.T) {
			rel := filepath_.GetParentRelPath(testCase.filePath)
			assert.Equal(t, testCase.expected, rel)
			t.Logf("rel: %v", rel)

		})
	}
}

func TestAbsPath(t *testing.T) {
	testCases := []struct {
		filePath string
		expected string
	}{
		{
			filePath: "",
			expected: "",
		},
		{
			filePath: ".",
			expected: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.filePath, func(t *testing.T) {
			abs, err := filepath.Abs(testCase.filePath)
			if err != nil {
				t.Fatalf("abs err: %v", err)
			}
			t.Logf("abs: %v", abs)

		})
	}
}

func TestJoinPaths(t *testing.T) {
	testCases := []struct {
		absolutePath string
		relativePath string
		expected     string
	}{
		{
			absolutePath: "/data",
			relativePath: "./relative",
			expected:     "",
		},
		{
			absolutePath: "./data",
			relativePath: "./relative",
			expected:     "",
		},
		{
			absolutePath: "./data",
			relativePath: "./relative/",
			expected:     "",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test-case-%d", i), func(t *testing.T) {
			fullPath, err := filepath_.JoinPaths(testCase.absolutePath, testCase.relativePath)
			if err != nil {
				t.Fatalf("joinPaths err: %v", err)
			}
			t.Logf("fullPath : %v", fullPath)
		})
	}
}
