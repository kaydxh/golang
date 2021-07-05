package net_test

import (
	"testing"

	net_ "github.com/kaydxh/golang/go/net"
)

func TestGetLocalAddrs(t *testing.T) {
	addrs, err := net_.GetLocalAddrs()
	if err != nil {
		t.Fatalf("failed to get local addrs, err: %v", err)
		return
	}
	t.Logf("addrs: %v", addrs)
}

func TestGetLocalIPs(t *testing.T) {
	ips, err := net_.GetLocalIPs()
	if err != nil {
		t.Fatalf("failed to get local addrs, err: %v", err)
		return
	}
	t.Logf("ips: %v", ips)
}

func TestGetLocalFirstIP(t *testing.T) {
	ip, err := net_.GetLocalFirstIP()
	if err != nil {
		t.Fatalf("failed to get local first addr, err: %v", err)
		return
	}
	t.Logf("ip: %v", ip)
}
