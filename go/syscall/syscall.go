package syscall

import (
	"syscall"

	errors_ "github.com/kaydxh/golang/go/errors"
)

func KillBatch(pids []int, sig syscall.Signal) (errorPids []int, err error) {
	var errs []error
	for _, pid := range pids {
		err = syscall.Kill(pid, sig)
		if err != nil {
			errorPids = append(errorPids, pid)
			errs = append(errs, err)
		}
	}

	return errorPids, errors_.NewAggregate(errs)
}
