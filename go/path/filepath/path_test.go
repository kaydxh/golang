package filepath_test

import (
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
