package k8sdns_test

import (
	"testing"

	k8sdns_ "github.com/kaydxh/golang/pkg/resolver/dns/k8s-resolver"
)

func TestK8sDNSResolverGetPodsUsingInformer(t *testing.T) {
	resolver, err := k8sdns_.NewK8sDNSResolver()
	if err != nil {
		t.Fatalf("new k8s dns reslover, err: %v", err)
	}

	pods, err := resolver.Pods()
	if err != nil {
		t.Fatalf("resolver get pods, err: %v", err)
	}

	t.Logf("get pods: %v", pods)
}
