package net_test

import (
	"net"
	"testing"

	net_ "github.com/kaydxh/golang/go/net"
	"github.com/stretchr/testify/assert"
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

func TestGetHostIP(t *testing.T) {
	ip, err := net_.GetHostIP()
	if err != nil {
		t.Fatalf("failed to get host ip, err: %v", err)
		return
	}
	t.Logf("ip: %v", ip)
}

func TestIsIPv4String(t *testing.T) {
	isIPv4 := net_.IsIPv4String("199.591.149.232")
	assert.Equal(t, true, isIPv4)
	t.Logf("ipv4: %v", isIPv4)
}

func TestLookupHost(t *testing.T) {
	ips, err := net.LookupHost("www.google.com")
	if err != nil {
		t.Fatalf("failed to get host ip, err: %v", err)
		return
	}
	t.Logf("ips: %v", ips)
}

func TestLookupHostIPv4(t *testing.T) {
	ips, err := net_.LookupHostIPv4("www.google.com")
	if err != nil {
		t.Fatalf("failed to get host ip, err: %v", err)
		return
	}
	t.Logf("ips: %v", ips)
}
