package resolver_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	resolver_ "github.com/kaydxh/golang/pkg/resolver"
	viper_ "github.com/kaydxh/golang/pkg/viper"
)

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
