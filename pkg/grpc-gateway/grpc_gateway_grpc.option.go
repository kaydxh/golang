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
package grpcgateway

import (
	"runtime/debug"
	"time"

	interceptorlogrus_ "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	interceptorrecovery_ "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	runtime_ "github.com/kaydxh/golang/go/runtime"
	rate_ "github.com/kaydxh/golang/go/time/rate"
	interceptortcloud_ "github.com/kaydxh/golang/pkg/middleware/api/tcloud/v3.0"
	interceptordebug_ "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/debug"
	interceptoropentelemetry_ "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/opentelemetry"
	interceptorratelimit_ "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/ratelimit"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QPSLimitConfig QPS限流配置
type QPSLimitConfig struct {
	// DefaultQPS 默认QPS限制，0表示不限制
	DefaultQPS float64
	// DefaultBurst 默认突发容量
	DefaultBurst int
	// MaxConcurrency 最大并发数限制，0表示不限制
	// 与QPS限流不同，并发控制限制的是同时处理的请求数，请求完成后令牌会归还
	MaxConcurrency int
	// MethodQPS 方法级QPS配置，key为完整方法名（如 /package.Service/Method）
	MethodQPS map[string]float64
	// MethodBurst 方法级突发容量配置
	MethodBurst map[string]int
	// MethodMaxConcurrency 方法级最大并发数配置
	MethodMaxConcurrency map[string]int
}

func WithServerOptions(opts ...grpc.ServerOption) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.serverOptions = append(c.opts.serverOptions, opts...)
	})
}

func WithClientDialOptions(opts ...grpc.DialOption) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.clientDialOptions = append(c.opts.clientDialOptions, opts...)
	})
}

func WithServerUnaryInterceptorsOptions(opts ...grpc.UnaryServerInterceptor) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.interceptionOptions.grpcServerOpts.unaryInterceptors = append(
			c.opts.interceptionOptions.grpcServerOpts.unaryInterceptors,
			opts...)
	})
}

func WithServerStreamInterceptorsOptions(opts ...grpc.StreamServerInterceptor) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.interceptionOptions.grpcServerOpts.streamInterceptors = append(
			c.opts.interceptionOptions.grpcServerOpts.streamInterceptors,
			opts...)
	})
}

func WithServerInterceptorsLogrusOptions(
	logger *logrus.Logger,
) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		l := logrus.NewEntry(logger)
		WithServerUnaryInterceptorsOptions(interceptorlogrus_.UnaryServerInterceptor(l)).apply(c)
		WithServerStreamInterceptorsOptions(interceptorlogrus_.StreamServerInterceptor(l)).apply(c)
	})
}

// recover
func WithServerInterceptorsRecoveryOptions() GRPCGatewayOption {
	opts := []interceptorrecovery_.Option{
		interceptorrecovery_.WithRecoveryHandler(func(r interface{}) (err error) {
			out, err := runtime_.FormatStack()
			if err != nil {
				logrus.WithError(status.Errorf(codes.Internal, "%s", r)).Errorf("%s", debug.Stack())
			} else {
				logrus.WithError(status.Errorf(codes.Internal, "%s", r)).Errorf("%s", out)
			}

			return status.Errorf(codes.Internal, "panic: %v", r)
		}),
	}

	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(interceptorrecovery_.UnaryServerInterceptor(opts...)).apply(c)
		WithServerStreamInterceptorsOptions(interceptorrecovery_.StreamServerInterceptor(opts...)).apply(c)
	})
}

// Deprecated: Use WithServerUnaryMetricInterceptorOptions instead
func WithServerUnaryInterceptorsTimerOptions(enabledMetric bool) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		if enabledMetric {
			WithServerUnaryInterceptorsOptions(interceptoropentelemetry_.UnaryServerMetricInterceptor()).apply(c)
		}
	})
}

// Deprecated: Use WithServerUnaryMetricInterceptorOptions instead
func WithServerUnaryInterceptorsCodeMessageOptions(enabledMetric bool) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		if enabledMetric {
			WithServerUnaryInterceptorsOptions(
				interceptoropentelemetry_.UnaryServerMetricInterceptor(),
			).apply(
				c,
			)
		}
	})
}

func WithServerUnaryMetricInterceptorOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(
			interceptoropentelemetry_.UnaryServerMetricInterceptor(),
		).apply(
			c,
		)
	})
}

// WithServerUnaryTraceInterceptorOptions adds OpenTelemetry trace interceptor
func WithServerUnaryTraceInterceptorOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(
			interceptoropentelemetry_.UnaryServerTraceInterceptor(),
		).apply(c)
		WithServerStreamInterceptorsOptions(
			interceptoropentelemetry_.StreamServerTraceInterceptor(),
		).apply(c)
	})
}

func WithServerUnaryInterceptorsRequestIdOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(interceptordebug_.UnaryServerInterceptorOfRequestId()).apply(c)
	})
}

func WithServerUnaryInterceptorsErrorOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(interceptortcloud_.UnaryServerInterceptorOfError()).apply(c)
		//WithServerStreamInterceptorsOptions(interceptortcloud_.StreamServerInterceptor()).apply(c)
	})
}

//limiter rate for grpc api
func WithServerInterceptorsLimitRateOptions(burstUnary, burstStream int) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		if burstUnary > 0 {
			limiterUnary := rate_.NewLimiter(burstUnary)
			WithServerUnaryInterceptorsOptions(interceptorratelimit_.UnaryServerInterceptor(limiterUnary)).apply(c)
		}

		if burstStream > 0 {
			limiterStream := rate_.NewLimiter(burstStream)
			WithServerStreamInterceptorsOptions(interceptorratelimit_.StreamServerInterceptor(limiterStream)).apply(c)
		}
	})
}

// WithServerInterceptorsQPSLimitOptions QPS限流拦截器选项
// 支持全局默认QPS和方法级QPS配置
func WithServerInterceptorsQPSLimitOptions(config QPSLimitConfig) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		// QPS限流
		if config.DefaultQPS > 0 || len(config.MethodQPS) > 0 {
			// 创建方法级QPS限流器
			defaultBurst := config.DefaultBurst
			if defaultBurst <= 0 {
				defaultBurst = int(config.DefaultQPS)
			}
			methodLimiter := rate_.NewMethodQPSLimiter(config.DefaultQPS, defaultBurst)

			// 设置方法级QPS配置
			for method, qps := range config.MethodQPS {
				burst := config.MethodBurst[method]
				if burst <= 0 {
					burst = int(qps)
				}
				methodLimiter.SetMethodQPS(method, qps, burst)
			}

			WithServerUnaryInterceptorsOptions(
				interceptorratelimit_.UnaryServerInterceptorQPS(methodLimiter),
			).apply(c)
		}

		// 并发控制
		if config.MaxConcurrency > 0 || len(config.MethodMaxConcurrency) > 0 {
			defaultConcurrency := config.MaxConcurrency
			if defaultConcurrency <= 0 {
				defaultConcurrency = 0 // 不限制
			}
			concurrencyLimiter := rate_.NewMethodLimiter(defaultConcurrency)

			// 设置方法级并发配置
			for method, maxConcurrency := range config.MethodMaxConcurrency {
				if maxConcurrency > 0 {
					concurrencyLimiter.AddLimiter(method, rate_.NewLimiter(maxConcurrency))
				}
			}

			WithServerUnaryInterceptorsOptions(
				interceptorratelimit_.UnaryServerInterceptorConcurrency(concurrencyLimiter),
			).apply(c)
		}
	})
}

func WithServerUnaryInterceptorsInOutPacketOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(interceptordebug_.UnaryServerInterceptorOfInOutputPrinter()).apply(c)
	})
}

// timeout
func WithServerInterceptorTimeoutOptions(timeout time.Duration) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		if timeout > 0 {
			runtime.DefaultContextTimeout = timeout
		}
	})
}

/*
// WithHttpHandlerInterceptorOptions
func WithHttpHandlerInterceptorOptions(opts ...http_.HandlerChainOption) GRPCGatewayOption {

	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.interceptionOptions.httpServerOpts.handlerChain.ApplyOptions(opts...)
	})
}

// WithHttpHandlerInterceptorOptions
func WithHttpHandlerInterceptorOptions(opts ...http_.HandlerInterceptorsOption) GRPCGatewayOption {

	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.interceptionOptions.httpServerOpts.httpInterceptors.ApplyOptions(opts...)
	})
}

func WithHttpHandlerInterceptor(handler func(http.Handler) http.Handler) GRPCGatewayOption {

	return WithHttpHandlerInterceptorOptions(
		http_.HandlerInterceptorsOptionFunc(func(handlers *http_.HandlerInterceptors) {
			handlers.Interceptors = append(handlers.Interceptors, handler)
		}),
	)
}
*/

// WithHttpHandlerInterceptorTraceIDOptions
/*
func WithHttpHandlerInterceptorInOutPacketOptions() GRPCGatewayOption {
	//return WithHttpHandlerInterceptor(httpinterceptorinoutpacket_.InOutPacket)
	return WithHttpHandlerInterceptor(httpinterceptorinoutpacket_.InOutPacket)
}
*/
