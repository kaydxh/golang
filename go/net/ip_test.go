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
	t.Logf("ip: %v", ip.String())
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

func TestGetServerName(t *testing.T) {
	serverName := net_.GetServerName()
	t.Logf("serverName: %v", serverName)
}
