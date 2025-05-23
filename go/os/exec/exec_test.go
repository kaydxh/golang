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
package exec_test

import (
	"fmt"
	"testing"
	"time"

	exec_ "github.com/kaydxh/golang/go/os/exec"
)

func TestExec(t *testing.T) {
	cmd := fmt.Sprintf(`ps -elf`)
	timeout := 5 * time.Second //ms
	result, msg, err := exec_.Exec(time.Duration(timeout), cmd)
	if err != nil {
		t.Errorf("expect nil, got %v, msg: %v", err, msg)
	}
	t.Logf("result: %v", result)
}

func TestExecTimeout(t *testing.T) {
	// sleep 2s
	cmd := fmt.Sprintf(`sleep 2`)
	timeout := time.Second
	result, msg, err := exec_.Exec(time.Duration(timeout), cmd)
	if err != nil {
		t.Errorf("expect nil, got %v, msg: %v", err, msg)
	}
	t.Logf("result: %v", result)
}
