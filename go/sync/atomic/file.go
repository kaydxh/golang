package atomic

import (
	"fmt"
	"os"
)

type File string

func (m *File) TryLock() error {
	if m == nil {
		return fmt.Errorf("nil pointer")
	}

	if *m == "" {
		return fmt.Errorf("invalid file")
	}

	f, err := os.OpenFile(string(*m), os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}

	return f.Close()
}

func (m *File) TryUnLock() error {
	if m == nil || *m == "" {
		return nil
	}

	return os.Remove(string(*m))
}
