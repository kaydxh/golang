package webserver

import (
	"time"

	gw_ "github.com/kaydxh/golang/pkg/grpc-gateway"
	"github.com/spf13/viper"
)

func WithViper(v *viper.Viper) ConfigOption {
	return ConfigOptionFunc(func(c *Config) {
		c.opts.viper = v
	})
}

func WithShutdownDelayDuration(shutdownDelayDuration time.Duration) ConfigOption {
	return ConfigOptionFunc(func(c *Config) {
		c.opts.shutdownDelayDuration = shutdownDelayDuration
	})
}

func WithGRPCGatewayOptions(opts ...gw_.GRPCGatewayOption) ConfigOption {
	return ConfigOptionFunc(func(c *Config) {
		c.opts.gatewayOptions = append(c.opts.gatewayOptions, opts...)
	})
}
