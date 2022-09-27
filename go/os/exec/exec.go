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
package exec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"

	context_ "github.com/kaydxh/golang/go/context"
)

type CommandBuilder struct {
	cmd exec.Cmd

	opts struct {
		Timeout time.Duration
	}
}

func NewCommandBuilder(
	options ...CommandBuilderOption,
) (*CommandBuilder, error) {
	c := &CommandBuilder{}
	c.ApplyOptions(options...)

	return c, nil
}

// Exec return output result, err messgae, error
func (c *CommandBuilder) Exec(cmdName string,
	args ...string,
) (string, string, error) {
	return Exec(c.opts.Timeout, cmdName, args...)
}

//timout ms
func Exec(
	timeout time.Duration,
	cmdName string,
	args ...string,
) (string, string, error) {
	ctx, cancel := context_.WithTimeout(
		context.Background(),
		timeout,
	)
	defer cancel()

	var stdout, stderr bytes.Buffer
	args = append([]string{"-c", cmdName}, args...)
	cmd := exec.CommandContext(ctx, "/bin/sh", args...)
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return "", "", fmt.Errorf(
			"error starting %v:\nCommand stdout:\n%v\nstderr:\n%v\nerror:\n%v",
			cmd,
			cmd.Stdout,
			cmd.Stderr,
			err,
		)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.Wait()
	}()
	select {
	case err := <-errCh:
		if err != nil {
			var rc = 127
			if ee, ok := err.(*exec.ExitError); ok {
				rc = int(ee.Sys().(syscall.WaitStatus).ExitStatus())
			}
			return stdout.String(), stderr.String(),
				fmt.Errorf(
					"error running %v:\nCommand stdout:\n%v\nstderr:\n%v\nerror:\n%v\ncode:\n%v",
					cmd,
					cmd.Stdout,
					cmd.Stderr,
					err,
					rc,
				)
		}
	case <-ctx.Done():
		cmd.Process.Kill()
		return "", "", fmt.Errorf(
			"timed out waiting for command %v:\nCommand stdout:\n%v\nstderr:\n%v",
			cmd,
			cmd.Stdout,
			cmd.Stderr,
		)
	}

	return stdout.String(), stderr.String(), nil
}
