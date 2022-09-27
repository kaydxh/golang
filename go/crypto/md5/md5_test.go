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
package md5_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	md5_ "github.com/kaydxh/golang/go/crypto/md5"
	"gotest.tools/assert"
)

func TestMd5File(t *testing.T) {
	file, err := ioutil.TempFile(".", "file")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	defer os.RemoveAll(file.Name())

	strContext := "hello world"
	//	n, err := io.WriteString(file, "hello world")
	n, err := file.Write([]byte(strContext))
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	assert.Equal(t, len(strContext), n)
	fmt.Printf("fileName: %v, n: %v\n", file.Name(), n)

	sum, err := md5_.SumFile(file.Name())
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	assert.Equal(t, sum, md5_.SumString(strContext))
}

func TestMd5FileAt(t *testing.T) {
	file, err := ioutil.TempFile(".", "file")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	defer os.RemoveAll(file.Name())

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

	var offset int64
	for i, testCase := range testCases {
		_, err := file.Write(testCase.words)
		if err != nil {
			t.Errorf("expect nil, got %v", err)
		}

		if i > 0 {
			offset += int64(len(testCases[i-1].words))
		}
		fmt.Printf("i: %d, offset: %v, testCase: %s\n", i, offset, testCase.words)
		sum, err := md5_.SumFileAt(file.Name(), offset, int64(len(testCase.words)))
		if err != nil {
			t.Errorf("expect nil, got %v", err)
		}
		fmt.Println("sum: ", sum)
		assert.Equal(t, sum, md5_.SumBytes(testCase.words))
	}

}
