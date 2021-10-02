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

func GetNumFiles() (uint64, error) {
	var (
		rlimit unix.Rlimit
		zero   unix.Rlimit
	)

	err := unix.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit)
	if err != nil {
		return 0, err
	}
	if rlimit == zero {
		return 0, fmt.Errorf("failed to get rlimit, got zero value: %#v", rlimit)
	}

	return rlimit.Cur, nil
}
