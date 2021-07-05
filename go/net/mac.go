package net

import (
	"fmt"
	"net"
)

func GetLocalFirstMac() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var macs []string
	for _, netInterface := range interfaces {
		mac := netInterface.HardwareAddr.String()
		if len(mac) == 0 {
			continue
		} else {
			macs = append(macs, mac)
		}
	}
	if len(macs) == 0 {
		return "", fmt.Errorf("no valid mac")
	}

	return macs[0], nil
}

func GetLocalMacs() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var macs []string
	for _, netInterface := range interfaces {
		mac := netInterface.HardwareAddr.String()
		if len(mac) == 0 {
			continue
		} else {
			macs = append(macs, mac)
		}
	}

	return macs, nil
}
