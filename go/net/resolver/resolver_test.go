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
