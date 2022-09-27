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
package syscall_test

import (
	"testing"

	syscall_ "github.com/kaydxh/golang/go/syscall"
)

func TestGetNumFiles(t *testing.T) {

	curOpenFiles, maxOpenFiles, err := syscall_.GetNumFiles()
	if err != nil {
		t.Fatalf("failed to get num files, err: %v", err)
	}

	t.Logf("curOpenFiles: %v,  maxOpenFiles: %v", curOpenFiles, maxOpenFiles)
}

func TestSetMaxNumFiles(t *testing.T) {

	curOpenFiles, err := syscall_.SetMaxNumFiles()
	if err != nil {
		t.Fatalf("failed to set max num files, err: %v", err)
	}

	t.Logf("curOpenFiles: %v", curOpenFiles)
}
