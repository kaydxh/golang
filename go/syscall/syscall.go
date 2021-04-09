package sysycall

import "syscall"

func KillBatch(pids []int, sig syscall.Signal) (errorPids []int, err error) {
	for _, pid := range pids {
		errIn := syscall.Kill(pid, sig)
		if errIn != nil {
			errorPids = append(errorPids, pid)
			err = errIn
		}
	}

	return errorPids, err
}
