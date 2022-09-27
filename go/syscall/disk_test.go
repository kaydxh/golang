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

// GOOS=linux  GOARCH=amd64  go test -c disk_test.go -o test
// /test -test.v
func TestDiskUsage(t *testing.T) {
	testCases := []struct {
		volumePath string
		expected   string
	}{
		{
			volumePath: "/dev",
			expected:   "",
		},
		{
			volumePath: "/data",
			expected:   "",
		},
		{
			volumePath: "/data/home/log",
			expected:   "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.volumePath, func(t *testing.T) {
			du, err := syscall_.NewDiskUsage(testCase.volumePath)
			if err != nil {
				t.Errorf("new disk for path[%v] err, got : %s", testCase.volumePath, err)
				return

			}
			t.Logf(
				"disk free[%v], avali[%v], size[%v], used[%v], usage: %v",
				du.Free(),
				du.Avail(),
				du.Size(),
				du.Used(),
				du.Usage(),
			)

		})
	}
}
