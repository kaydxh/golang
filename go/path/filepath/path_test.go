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
