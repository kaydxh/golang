package k8sdns_test

import (
	"context"
	"flag"
	"testing"

	k8sdns_ "github.com/kaydxh/golang/pkg/resolver/dns/k8s-resolver"
)

var svc string

func init() {
	flag.StringVar(&svc, "s", "test-cube-algo-backend", "k8s service")
}

func TestK8sDNSResolverGetPodsUsingInformer(t *testing.T) {
	resolver, err := k8sdns_.NewK8sDNSResolver()
	if err != nil {
		t.Fatalf("new k8s dns reslover, err: %v", err)
	}

	t.Logf("svc: %v", svc)
	pods, err := resolver.Pods(context.Background(), svc)
	if err != nil {
		t.Fatalf("resolver get pods, err: %v", err)
	}

	t.Logf("get pods: %v", pods)
}
