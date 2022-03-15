package reslover_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kaydxh/golang/pkg/reslover"
	reslover_ "github.com/kaydxh/golang/pkg/reslover"
	viper_ "github.com/kaydxh/golang/pkg/viper"
)

func TestNewResloverService(t *testing.T) {
	cfgFile := "./reslover.yaml"
	config := reslover_.NewConfig(reslover_.WithViper(viper_.GetViper(cfgFile, "reslover")))
	s, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}
	s.Run(context.Background())

	type args struct {
		consistkey       string
		resolverInterval time.Duration
		services         []reslover_.ResloverQuery
	}
	tests := []struct {
		name string
		args args
		want *reslover_.ResloverService
	}{
		// TODO: Add test cases.
		{
			name: "www.baidu.com",
			args: args{
				consistkey:       "1",
				resolverInterval: 0,
				services: []reslover_.ResloverQuery{
					{
						Domain: "www.baidu.com",
						Opts: reslover_.ResloverOptions{
							ResloverType:    reslover_.Reslover_reslover_type_dns,
							LoadBalanceMode: reslover.Reslover_load_balance_mode_consist,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := reslover_.NewResloverService(tt.args.resolverInterval, tt.args.services...)
			for i := 0; i < 100; i++ {
				consistkey := fmt.Sprintf("consist-key-%d", i)
				node, has := rs.PickNode(tt.name, consistkey)
				t.Logf("pick node: %v, has: %v, consistkey: %v", node, has, consistkey)
			}
		})
	}
}
