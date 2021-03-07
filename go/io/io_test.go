package io_test

import (
	"bytes"
	"strings"
	"testing"

	io_ "github.com/kaydxh/golang/go/io"
)

func TestWriteLine(t *testing.T) {
	testCases := []struct {
		name     string
		words    []string
		expected string
	}{
		{
			name:     "write no word",
			words:    []string{},
			expected: "",
		},
		{
			name:     "write one word",
			words:    []string{"test1"},
			expected: "test1\n",
		},

		{
			name:     "write multi words",
			words:    []string{"test1", "test2", "test3"},
			expected: "test1 test2 test3\n",
		},
	}

	testBuffer := bytes.NewBuffer(nil)
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testBuffer.Reset()
			io_.WriteLine(testBuffer, testCase.words...)
			if !strings.EqualFold(testBuffer.String(), testCase.expected) {
				t.Fatalf(
					"write word is %v\n expected: %s, got: %s",
					testCase.words,
					testCase.expected,
					testBuffer.String(),
				)

			}
		})
	}
}
