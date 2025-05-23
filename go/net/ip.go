/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
package net

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	hash_ "github.com/kaydxh/golang/go/encoding/hash"
	os_ "github.com/kaydxh/golang/go/os"
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
	return LookupHostIPv4WithContext(context.Background(), host)
}

func LookupHostIPv4WithContext(ctx context.Context, host string) (addrs []string, err error) {
	ips, err := net.DefaultResolver.LookupHost(ctx, host)
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

// SplitHostIntPort split host and integral port
func SplitHostIntPort(s string) (string, int, error) {
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return "", 0, err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return "", 0, err
	}
	return host, portInt, err
}

// parseTarget takes the user input target string and default port, returns formatted host and port info.
// If target doesn't specify a port, set the port to be the defaultPort.
// If target is in IPv6 format and host-name is enclosed in square brackets, brackets
// are stripped when setting the host.
// examples:
// target: "www.google.com" defaultPort: "443" returns host: "www.google.com", port: "443"
// target: "ipv4-host:80" defaultPort: "443" returns host: "ipv4-host", port: "80"
// target: "[ipv6-host]" defaultPort: "443" returns host: "ipv6-host", port: "443"
// target: ":80" defaultPort: "443" returns host: "localhost", port: "80"
func ParseTarget(target, defaultPort string) (host, port string, err error) {
	if target == "" {
		return "", "", fmt.Errorf("target is empty")
	}
	if ip := net.ParseIP(target); ip != nil {
		// target is an IPv4 or IPv6(without brackets) address
		return target, defaultPort, nil
	}
	if host, port, err = net.SplitHostPort(target); err == nil {
		if port == "" {
			// If the port field is empty (target ends with colon), e.g. "[::1]:", this is an error.
			return "", "", fmt.Errorf("missing port after port-separator colon")
		}
		// target has port, i.e ipv4-host:port, [ipv6-host]:port, host-name:port
		if host == "" {
			// Keep consistent with net.Dial(): If the host is empty, as in ":80", the local system is assumed.
			host = "localhost"
		}
		return host, port, nil
	}
	if host, port, err = net.SplitHostPort(target + ":" + defaultPort); err == nil {
		// target doesn't have port
		return host, port, nil
	}
	return "", "", fmt.Errorf("invalid target address %v, error info: %v", target, err)
}

// 1 get k8s pod name
// 2 get ip:pid
func GetServerName() string {
	serverName := fmt.Sprintf("%v:%v", os.Getenv("POD_NAMESPACE"), os.Getenv("POD_NAME"))
	if serverName == ":" {
		ip, _ := GetHostIP()
		pid := os_.GetProcId()
		return fmt.Sprintf("%s:%v", ip.String(), pid)
	}
	return serverName
}

func GetServerId() uint32 {
	return hash_.HashCode(GetServerName())
}
