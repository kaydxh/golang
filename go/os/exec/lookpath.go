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
package exec

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// LookPathIn resolves a bare command name against the PATH entry of environ,
// the environment a child process will actually get, instead of the current
// process's PATH. Relative PATH entries are skipped, as the child's working
// directory would decide them. A missing name wraps exec.ErrNotFound.
func LookPathIn(file string, environ []string) (string, error) {
	if strings.ContainsAny(file, `/\`) {
		return "", fmt.Errorf("lookpath %q: not a bare command name", file)
	}
	for i := len(environ) - 1; i >= 0; i-- {
		name, value, ok := strings.Cut(environ[i], "=")
		if !ok || name != "PATH" {
			continue
		}
		for _, dir := range filepath.SplitList(value) {
			if dir == "" || !filepath.IsAbs(dir) {
				continue
			}
			if path, err := exec.LookPath(filepath.Join(dir, file)); err == nil {
				return path, nil
			}
		}
		break
	}
	return "", &exec.Error{Name: file, Err: exec.ErrNotFound}
}
