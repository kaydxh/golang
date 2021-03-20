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
