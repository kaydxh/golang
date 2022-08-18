package resolver_test

import (
	"fmt"
	"testing"

	"context"

	resolver_ "github.com/kaydxh/golang/go/net/resolver"
	_ "github.com/kaydxh/golang/go/net/resolver/dns"
)

func TestResolveOne(t *testing.T) {
	testCases := []struct {
		target   string
		expected string
	}{
		{
			target:   "dns:///www.baidu.com",
			expected: "",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d-test", i), func(t *testing.T) {
			addr, err := resolver_.ResolveOne(context.Background(), testCase.target)
			if err != nil {
				t.Fatalf("failed to resolve target: %v, err: %v", testCase.target, err)
			}
			t.Logf("resolve one addr %v for target %v", addr, testCase.target)
		})
	}
}
