package io_test

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	io_ "github.com/kaydxh/golang/go/io"
	os_ "github.com/kaydxh/golang/go/os"
	"github.com/stretchr/testify/assert"
)

func TestWriteReadLine(t *testing.T) {
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

	lines := io_.ReadLines(testBuffer.Bytes())
	assert.Equal(t, "test1 test2 test3", lines[0])

}

func TestWriteReadFileLine(t *testing.T) {
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

	workDir, _ := os_.Getwd()
	testFileNew := filepath.Join(workDir, "test-file-new")
	testFileAppend := filepath.Join(workDir, "test-file-append")

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := io_.WriteFileLine(testFileNew, testCase.words, false)
			assert.Nil(t, err)
			err = io_.WriteFileLine(testFileAppend, testCase.words, true)
			assert.Nil(t, err)
			/*
				if !strings.EqualFold(testBuffer.String(), testCase.expected) {
					t.Fatalf(
						"write word is %v\n expected: %s, got: %s",
						testCase.words,
						testCase.expected,
						testBuffer.String(),
					)

				}
			*/
		})
	}

	lines, err := io_.ReadFileLines(testFileAppend)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	fmt.Println("read lines:", lines)
}

func TestWriteFileLines(t *testing.T) {
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

	workDir, _ := os_.Getwd()
	testFileNew := filepath.Join(workDir, "test-file-new")
	testFileAppend := filepath.Join(workDir, "test-file-append")

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := io_.WriteFileLines(testFileNew, testCase.words, false)
			assert.Nil(t, err)
			err = io_.WriteFileLines(testFileAppend, testCase.words, true)
			assert.Nil(t, err)
			/*
				if !strings.EqualFold(testBuffer.String(), testCase.expected) {
					t.Fatalf(
						"write word is %v\n expected: %s, got: %s",
						testCase.words,
						testCase.expected,
						testBuffer.String(),
					)

				}
			*/
		})
	}
}
