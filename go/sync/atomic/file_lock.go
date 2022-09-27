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
package atomic

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	os_ "github.com/kaydxh/golang/go/os"
)

type FileLock string

func (m *FileLock) TryLock() error {
	if m == nil {
		return fmt.Errorf("nil pointer")
	}
	if *m == "" {
		return fmt.Errorf("invalid file")
	}
	name := string(*m)
	if !strings.HasSuffix(name, ".lock") {
		*m += ".lock"
	}

	dir := filepath.Dir(name)
	err := os_.MakeDirAll(dir)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(string(*m), os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}

	return f.Close()
}

func (m *FileLock) TryUnLock() error {
	if m == nil || *m == "" {
		return nil
	}

	return os.Remove(string(*m))
}

/*
func makePidFile(name string) (tmpname string, cleanup func(), err error) {
	tmplock, err := ioutil.TempFile(
		filepath.Dir(name),
		filepath.Base(name)+fmt.Sprintf("%d", os.GetPid())+".lock",
	)
	if err != nil {
		return "", nil, err
	}

	cleanup = func() {
		_ = tmplock.Close()
		_ = os.Remove(tmplock.Name())
	}

}
*/
