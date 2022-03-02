package net

import (
	"errors"
	"fmt"
	"net"
)

// This code is borrowed from https://github.com/uber/tchannel-go/blob/dev/localip.go

// scoreAddr scores how likely the given addr is to be a remote address and returns the
// IP to use when listening. Any address which receives a negative score should not be used.
// Scores are calculated as:
// -1 for any unknown IP addresses.
// +300 for IPv4 addresses
// +100 for non-local addresses, extra +100 for "up" interaces.
func scoreAddr(iface net.Interface, addr net.Addr) (int, net.IP) {
	var ip net.IP
	if netAddr, ok := addr.(*net.IPNet); ok {
		ip = netAddr.IP
	} else if netIP, ok := addr.(*net.IPAddr); ok {
		ip = netIP.IP
	} else {
		return -1, nil
	}

	var score int
	if ip.To4() != nil {
		score += 300
	}
	if iface.Flags&net.FlagLoopback == 0 && !ip.IsLoopback() {
		score += 100
		if iface.Flags&net.FlagUp != 0 {
			score += 100
		}
	}

	fmt.Printf("ip: %v, score: %v\n", ip, score)
	return score, ip
}

// HostIP tries to find an IP that can be used by other machines to reach this machine.
func GetHostIP() (net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	bestScore := -1
	var bestIP net.IP
	// Select the highest scoring IP as the best IP.
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			// Skip this interface if there is an error.
			continue
		}

		for _, addr := range addrs {
			score, ip := scoreAddr(iface, addr)
			if score > bestScore {
				bestScore = score
				bestIP = ip
			}
		}
	}

	if bestScore == -1 {
		return nil, errors.New("no addresses to listen on")
	}

	return bestIP, nil
}

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
		ipNet, ok := addr.(*net.IPNet)
		if ok {
			localAddrs = append(localAddrs, ipNet)
		}
	}

	return localAddrs, nil
}

// IsIPv4 returns if netIP is IPv4.
func IsIPv4(netIP net.IP) bool {
	return netIP != nil && netIP.To4() != nil
}

func IsIPv4String(ip string) bool {
	netIP := ParseIP(ip)
	return IsIPv4(netIP)
}

func LookupHostIPv4(host string) (addrs []string, err error) {
	ips, err := net.LookupHost(host)
	if err != nil {
		return
	}

	for _, ip := range ips {
		if IsIPv4String(ip) {
			addrs = append(addrs, ip)
		}
	}

	return addrs, err
}
