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
//go:build windows

package syscall

import (
	"syscall"

	"golang.org/x/sys/windows"

	errors_ "github.com/kaydxh/golang/go/errors"
)

// KillBatch terminates processes by pid. Windows has no POSIX signals: the
// only cross-process signal is termination, so every call terminates
// regardless of the requested signal (the signature is preserved for the
// shared callers; the sig value is ignored).
func KillBatch(pids []int, _ syscall.Signal) (errorPids []int, err error) {
	var errs []error
	for _, pid := range pids {
		h, err := windows.OpenProcess(windows.PROCESS_TERMINATE, false, uint32(pid))
		if err != nil {
			errorPids = append(errorPids, pid)
			errs = append(errs, err)
			continue
		}
		if err := windows.TerminateProcess(h, 1); err != nil {
			errorPids = append(errorPids, pid)
			errs = append(errs, err)
		}
		windows.CloseHandle(h)
	}
	return errorPids, errors_.NewAggregate(errs)
}
