//go:build !windows

/*
 *Copyright (c) 2026, kaydxh
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
	"context"
	"os/exec"
	"testing"

	exec_ "github.com/kaydxh/golang/go/os/exec"
)

func TestSetProcessGroup(t *testing.T) {
	t.Parallel()
	cmd := exec.Command("echo")
	exec_.SetProcessGroup(cmd)
	if cmd.SysProcAttr == nil || !cmd.SysProcAttr.Setpgid {
		t.Fatalf("SysProcAttr = %+v", cmd.SysProcAttr)
	}
}

func TestKillProcessGroupNil(t *testing.T) {
	t.Parallel()
	if err := exec_.KillProcessGroup(nil); err != nil {
		t.Fatalf("err = %v", err)
	}
}

func TestCommandContextUsesProcessGroup(t *testing.T) {
	t.Parallel()
	cmd := exec_.CommandContext(context.Background(), "echo")
	if cmd.SysProcAttr == nil || !cmd.SysProcAttr.Setpgid || cmd.Cancel == nil {
		t.Fatalf("cmd = %+v", cmd)
	}
}
