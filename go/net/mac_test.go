package net_test

import (
	"testing"

	net_ "github.com/kaydxh/golang/go/net"
)

func TestGetLocalMacs(t *testing.T) {
	macs, err := net_.GetLocalMacs()
	if err != nil {
		t.Fatalf("failed to get local addrs, err: %v", err)
		return
	}
	t.Logf("macs: %v", macs)
}

func TestGetLocalFirstMac(t *testing.T) {
	mac, err := net_.GetLocalFirstMac()
	if err != nil {
		t.Fatalf("failed to get local first addr, err: %v", err)
		return
	}
	t.Logf("mac: %v", mac)
}
