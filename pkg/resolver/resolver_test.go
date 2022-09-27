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
	"context"
	"fmt"
	"testing"
	"time"

	resolver_ "github.com/kaydxh/golang/pkg/resolver"
	viper_ "github.com/kaydxh/golang/pkg/viper"
)

/*
func TestNewResolverService(t *testing.T) {
	cfgFile := "./resolver.yaml"
	config := resolver_.NewConfig(resolver_.WithViper(viper_.GetViper(cfgFile, "resolver")))
	s, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}
	s.Run(context.Background())
	//	net.DefaultResolver

	type args struct {
		consistkey       string
		resolverInterval time.Duration
		services         []resolver_.ResolverQuery
	}
	tests := []struct {
		name string
		args args
		want *resolver_.ResolverService
	}{
		// TODO: Add test cases.
		{
			name: "www.baidu.com",
			args: args{
				consistkey:       "1",
				resolverInterval: 0,
				services: []resolver_.ResolverQuery{
					{
						Domain: "www.baidu.com",
						Opts: resolver_.ResolverOptions{
							ResolverType:    resolver_.Resolver_resolver_type_dns,
							LoadBalanceMode: resolver_.Resolver_load_balance_mode_consist,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = s.AddServices(tt.args.services...)
			for i := 0; i < 100; i++ {
				consistkey := fmt.Sprintf("consist-key-%d", i)
				node, has := s.PickNode(tt.name, consistkey)
				t.Logf("pick node: %v, has: %v, consistkey: %v", node, has, consistkey)
			}
		})
	}
}
*/

func TestNewResolverService2(t *testing.T) {
	cfgFile := "./resolver.yaml"
	config := resolver_.NewConfig(resolver_.WithViper(viper_.GetViper(cfgFile, "resolver")))
	s, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}
	s.Run(context.Background())

	type args struct {
		consistkey       string
		resolverInterval time.Duration
		services         []resolver_.ResolverQuery
	}

	tests := []struct {
		serviceName      string
		nodeGroup        string
		nodeUnit         string
		consistkey       string
		loadBalanceMode  resolver_.Resolver_LoadBalanceMode
		resolverType     resolver_.Resolver_ResolverType
		resolverInterval time.Duration
	}{
		// TODO: Add test cases.
		{
			serviceName:      "test-cube-algo-backend",
			nodeGroup:        "edge-global-group",
			nodeUnit:         "edge-node-zone",
			consistkey:       "1",
			resolverInterval: 0,
			resolverType:     resolver_.Resolver_resolver_type_k8s,
		},
	}
	for _, tt := range tests {
		t.Run(tt.serviceName, func(t *testing.T) {
			rq, err := resolver_.NewResolverQuery(
				tt.serviceName,
				resolver_.WithLoadBalanceMode(tt.loadBalanceMode),
				resolver_.WithResolverType(tt.resolverType),
				resolver_.WithNodeGroup(tt.nodeGroup),
				resolver_.WithNodeUnit(tt.nodeUnit),
			)
			if err != nil {
				t.Fatalf("new resolver query err: %v", err)
			}
			err = s.AddServices(rq)
			time.Sleep(5 * time.Second)
			for i := 0; i < 100; i++ {
				consistkey := fmt.Sprintf("consist-key-%d", i)
				node, has := s.PickNode(tt.serviceName, consistkey)
				t.Logf("pick node: %v, has: %v, consistkey: %v", node, has, consistkey)
			}
		})
	}
}
