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
package syscall

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/unix"
)

// SetNumFiles sets the linux rlimit for the maximum open files.
func SetNumFiles(maxOpenFiles uint64) error {
	return unix.Setrlimit(unix.RLIMIT_NOFILE, &unix.Rlimit{Max: maxOpenFiles, Cur: maxOpenFiles})
}

func GetNumFiles() (uint64, uint64, error) {
	var (
		rlimit unix.Rlimit
		zero   unix.Rlimit
	)

	err := unix.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit)
	if err != nil {
		return 0, 0, err
	}
	if rlimit == zero {
		return 0, 0, fmt.Errorf("failed to get rlimit, got zero value: %#v", rlimit)
	}

	return rlimit.Cur, rlimit.Max, nil
}

func SetMaxNumFiles() (uint64, error) {

	_, maxOpenFiles, err := GetNumFiles()
	if err != nil {
		return 0, err
	}

	err = SetNumFiles(maxOpenFiles)
	if err != nil {
		return 0, err
	}

	newCurOpenFiles, _, err := GetNumFiles()
	if err != nil {
		return 0, err
	}
	if newCurOpenFiles != maxOpenFiles {
		return 0, fmt.Errorf("failed to set %d files, current open %v files", maxOpenFiles, newCurOpenFiles)

	}

	return newCurOpenFiles, nil
}
