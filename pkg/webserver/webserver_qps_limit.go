/*
 *Copyright (c) 2024, kaydxh
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
	gw_ "github.com/kaydxh/golang/pkg/grpc-gateway"
)

// QPSLimitConfig QPS限流配置
type QPSLimitConfig struct {
	// DefaultQPS 默认QPS限制，0表示不限制
	DefaultQPS float64 `json:"default_qps" yaml:"default_qps"`
	// DefaultBurst 默认突发容量，0表示使用DefaultQPS值
	DefaultBurst int `json:"default_burst" yaml:"default_burst"`
	// MaxConcurrency 最大并发数限制，0表示不限制
	// 与QPS限流不同，并发控制限制的是同时处理的请求数，请求完成后令牌会归还
	MaxConcurrency int `json:"max_concurrency" yaml:"max_concurrency"`
	// MethodQPS 方法/路径级QPS配置
	MethodQPS []MethodQPSConfigItem `json:"method_qps" yaml:"method_qps"`
}

// MethodQPSConfigItem 方法级QPS配置项
type MethodQPSConfigItem struct {
	// Method 方法名（gRPC: /package.Service/Method）或路径（HTTP: /api/v1/users）
	Method string `json:"method" yaml:"method"`
	// QPS 每秒请求数限制
	QPS float64 `json:"qps" yaml:"qps"`
	// Burst 突发容量，0表示使用QPS值
	Burst int `json:"burst" yaml:"burst"`
	// MaxConcurrency 方法级最大并发数限制，0表示不限制
	MaxConcurrency int `json:"max_concurrency" yaml:"max_concurrency"`
}

// ExtendedConfig 扩展配置，用于支持proto中未定义的字段
// 通过yaml配置文件或代码方式设置
type ExtendedConfig struct {
	// GrpcQPSLimit gRPC QPS限流配置
	GrpcQPSLimit *QPSLimitConfig `json:"grpc_qps_limit" yaml:"grpc_qps_limit"`
	// HttpQPSLimit HTTP QPS限流配置
	HttpQPSLimit *QPSLimitConfig `json:"http_qps_limit" yaml:"http_qps_limit"`
}

// ToGRPCQPSLimitConfig 转换为gRPC网关QPS限流配置
func (c *QPSLimitConfig) ToGRPCQPSLimitConfig() gw_.QPSLimitConfig {
	config := gw_.QPSLimitConfig{
		DefaultQPS:           c.DefaultQPS,
		DefaultBurst:         c.DefaultBurst,
		MaxConcurrency:       c.MaxConcurrency,
		MethodQPS:            make(map[string]float64),
		MethodBurst:          make(map[string]int),
		MethodMaxConcurrency: make(map[string]int),
	}
	for _, item := range c.MethodQPS {
		config.MethodQPS[item.Method] = item.QPS
		if item.Burst > 0 {
			config.MethodBurst[item.Method] = item.Burst
		}
		if item.MaxConcurrency > 0 {
			config.MethodMaxConcurrency[item.Method] = item.MaxConcurrency
		}
	}
	return config
}

// ToHTTPQPSLimitConfig 转换为HTTP QPS限流配置
func (c *QPSLimitConfig) ToHTTPQPSLimitConfig() gw_.HTTPQPSLimitConfig {
	config := gw_.HTTPQPSLimitConfig{
		DefaultQPS:         c.DefaultQPS,
		DefaultBurst:       c.DefaultBurst,
		MaxConcurrency:     c.MaxConcurrency,
		PathQPS:            make(map[string]float64),
		PathBurst:          make(map[string]int),
		PathMaxConcurrency: make(map[string]int),
	}
	for _, item := range c.MethodQPS {
		config.PathQPS[item.Method] = item.QPS
		if item.Burst > 0 {
			config.PathBurst[item.Method] = item.Burst
		}
		if item.MaxConcurrency > 0 {
			config.PathMaxConcurrency[item.Method] = item.MaxConcurrency
		}
	}
	return config
}
