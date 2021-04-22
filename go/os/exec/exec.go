package exec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"
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

//timout ms
func Exec(
	timeout time.Duration,
	name string,
	args ...string,
) (string, string, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(timeout)*time.Millisecond,
	)
	defer cancel()

	var stdout, stderr bytes.Buffer
	args = append([]string{"-c", name}, args...)
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
