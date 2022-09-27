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
