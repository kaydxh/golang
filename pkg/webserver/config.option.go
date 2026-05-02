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
package webserver

import (
	"time"

	gw_ "github.com/kaydxh/golang/pkg/grpc-gateway"
	healthz_ "github.com/kaydxh/golang/pkg/webserver/controller/healthz"
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

// WithGRPCQPSLimit 设置gRPC QPS限流配置
func WithGRPCQPSLimit(config *QPSLimitConfig) ConfigOption {
	return ConfigOptionFunc(func(c *Config) {
		c.opts.grpcQPSLimit = config
	})
}

// WithHTTPQPSLimit 设置HTTP QPS限流配置
func WithHTTPQPSLimit(config *QPSLimitConfig) ConfigOption {
	return ConfigOptionFunc(func(c *Config) {
		c.opts.httpQPSLimit = config
	})
}

// WithQPSLimit 同时设置gRPC和HTTP QPS限流配置
func WithQPSLimit(grpcConfig, httpConfig *QPSLimitConfig) ConfigOption {
	return ConfigOptionFunc(func(c *Config) {
		c.opts.grpcQPSLimit = grpcConfig
		c.opts.httpQPSLimit = httpConfig
	})
}

// WithHealthzOptions 设置 healthz 控制器的选项。
// 例如：WithHealthzOptions(healthz.WithDisableRootRoute()) 可禁用根路径健康检查，
// 避免与 SPA 前端的 static 控制器冲突。
func WithHealthzOptions(opts ...healthz_.ControllerOption) ConfigOption {
	return ConfigOptionFunc(func(c *Config) {
		c.opts.healthzOptions = opts
	})
}
