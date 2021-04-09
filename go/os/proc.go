package os

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func GetPidByName(
	name string,
) ([]int, error) {
	cmd := fmt.Sprintf(
		`ps ux | grep '%s'| grep -v grep | awk '{print $2}'`,
		name,
	)
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return nil, err
	}

	pids := strings.TrimSpace(string(result))
	sPids := strings.Split(pids, "\n")
	nPids := make([]int, 0)
	for _, pid := range sPids {
		nPid, err := strconv.Atoi(pid)
		if err != nil {
			return nil, err
		}
		nPids = append(nPids, nPid)
	}

	return nPids, nil
}
