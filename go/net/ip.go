package net

import (
	"fmt"
	"net"
)

func GetLocalFirstIP() (string, error) {
	ips, err := GetLocalIPs()
	if err != nil {
		return "", err
	}
	if len(ips) == 0 {
		return "", fmt.Errorf("no valid ip")
	}

	return ips[0], nil
}

func GetLocalIPs() ([]string, error) {
	localAddrs, err := GetLocalAddrs()
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, addr := range localAddrs {
		if addr.IP.To4() != nil && !addr.IP.IsLoopback() {
			ips = append(ips, addr.IP.String())
		}

	}

	return ips, nil
}

func GetLocalAddrs() ([]*net.IPNet, error) {
	var localAddrs []*net.IPNet

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		_, ipNet, err := net.ParseCIDR(addr.String())
		if err != nil {
			return nil, err
		}

		localAddrs = append(localAddrs, ipNet)
	}

	return localAddrs, nil
}
