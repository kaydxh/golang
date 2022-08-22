package resolver_test

import (
	"fmt"
	"testing"

	"context"

	"github.com/kaydxh/golang/go/net/resolver"
	resolver_ "github.com/kaydxh/golang/go/net/resolver"
	_ "github.com/kaydxh/golang/go/net/resolver/dns"
	resolve_ "github.com/kaydxh/golang/go/net/resolver/resolve"
)

func TestResolveOne(t *testing.T) {
	testCases := []struct {
		target   string
		iptype   resolver.Resolver_IPType
		expected string
	}{
		{
			target:   "dns:///www.baidu.com",
			iptype:   resolver.Resolver_ip_type_v4,
			expected: "",
		},
		{
			target:   "dns:///www.google.com",
			iptype:   resolver.Resolver_ip_type_v4,
			expected: "",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d-test", i), func(t *testing.T) {
			addr, err := resolve_.ResolveOne(context.Background(), testCase.target, resolver_.WithPickMode(resolver_.Resolver_pick_mode_first))
			if err != nil {
				t.Fatalf("failed to resolve target: %v, err: %v", testCase.target, err)
			}
			t.Logf("resolve one addr %v for target %v", addr, testCase.target)
		})
	}
}

func TestResolveAll(t *testing.T) {
	testCases := []struct {
		target   string
		iptype   resolver.Resolver_IPType
		expected string
	}{
		{
			target:   "dns:///www.baidu.com",
			iptype:   resolver.Resolver_ip_type_v4,
			expected: "",
		},
		{
			target:   "dns:///www.google.com",
			iptype:   resolver.Resolver_ip_type_v4,
			expected: "",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d-test", i), func(t *testing.T) {
			addrs, err := resolve_.ResolveAll(context.Background(), testCase.target, resolver.WithIPTypeForResolveAll(testCase.iptype))
			if err != nil {
				t.Fatalf("failed to resolve target: %v, err: %v", testCase.target, err)
			}
			t.Logf("resolve all addrs %v for target %v", addrs, testCase.target)
		})
	}
}
