package os

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	exec_ "github.com/kaydxh/golang/go/os/exec"
)

func GetPidsByName(
	timeout time.Duration,
	name string) ([]int, string, error) {
	cmd := fmt.Sprintf(
		`ps ux | grep '%s'| grep -v grep | awk '{print $2}'`,
		name,
	)
	result, msg, err := exec_.Exec(timeout, cmd)
	if err != nil {
		return nil, msg, err
	}

	pids := strings.TrimSpace(string(result))
	sPids := strings.Split(pids, "\n")
	var nPids []int
	for _, pid := range sPids {
		nPid, err := strconv.Atoi(pid)
		if err != nil {
			return nil, msg, err
		}
		nPids = append(nPids, nPid)
	}

	return nPids, msg, nil
}
