package io_test

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	md5_ "github.com/kaydxh/golang/go/crypto/md5"
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

func TestWriteAtOneThread(t *testing.T) {
	testCases := []struct {
		name     string
		words    []byte
		expected string
	}{
		{
			name:     "write one word",
			words:    []byte("test1"),
			expected: "test1",
		},
		{
			name:     "write one word",
			words:    []byte("test2"),
			expected: "test2",
		},

		{
			name:     "write one word",
			words:    []byte("test3"),
			expected: "test3",
		},
	}

	var offsets []int64 = []int64{10, 0, 5}
	workDir, _ := os_.Getwd()
	testFileOffset := filepath.Join(workDir, "test-file-offset")
	for i, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var offset int64
			if i > 0 {
				offset += int64(i * len(testCases[i-1].words))
			}
			fmt.Println("offset: ", offsets[i])
			//err := io_.WriteBytesAt(testFileOffset, testCases[i].words, offset)
			err := io_.WriteBytesAt(testFileOffset, testCases[i].words, offsets[i])
			assert.Nil(t, err)
		})
	}
}

func TestWriteAtMutilThreads(t *testing.T) {
	testCases := []struct {
		name     string
		words    []byte
		expected string
	}{
		{
			name:     "write one word",
			words:    []byte("test1"),
			expected: "test1",
		},
		{
			name:     "write one word",
			words:    []byte("test2"),
			expected: "test2",
		},

		{
			name:     "write one word",
			words:    []byte("test3"),
			expected: "test3",
		},
	}

	workDir, _ := os_.Getwd()
	testFileOffset := filepath.Join(workDir, "test-file-offset")
	wg := sync.WaitGroup{}
	for i, testCase := range testCases {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			t.Run(testCase.name, func(t *testing.T) {
				var offset int64
				for j := 0; j < index; j++ {
					offset += int64(len(testCases[j].words))
				}
				fmt.Printf("write: %s, offset: %d", testCases[index].words, offset)
				err := io_.WriteBytesAt(testFileOffset, testCases[index].words, offset)
				assert.Nil(t, err)
			})
		}(i)
	}
	wg.Wait()
}

func TestWriteReaderAtOneThread(t *testing.T) {
	testCases := []struct {
		name     string
		words    []string
		expected string
	}{
		{
			name:     "write one word",
			words:    []string{"test1"},
			expected: "test1",
		},
		{
			name:     "write one word",
			words:    []string{"test2"},
			expected: "test2",
		},

		{
			name:     "write one word",
			words:    []string{"test3"},
			expected: "test3",
		},
	}

	workDir, _ := os_.Getwd()
	testFile := filepath.Join(workDir, "test-file")
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := io_.WriteFileLine(testFile, testCase.words, true)
			assert.Nil(t, err)
		})
	}

	file, err := os_.OpenFile(testFile, true)
	if err != nil {
		return
	}
	defer file.Close()

	/*
		fileInfo, err := file.Stat()
		if err != nil {
			return
		}
	*/

	testFileCopy := filepath.Join(workDir, "test-file-copy")
	var offset int64

	for _, testCase := range testCases {

		err = io_.WriteReaderAt(testFileCopy, file, offset, int64(len(testCase.words[0])+1))
		offset += int64(len(testCase.words[0]) + 1)
		assert.Nil(t, err)
	}

	sumTestFile, _ := md5_.SumFile(testFile)
	sumTestFileCopy, _ := md5_.SumFile(testFileCopy)

	assert.Equal(t, sumTestFile, sumTestFileCopy)

}
