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
	"net/http"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	http_ "github.com/kaydxh/golang/go/net/http"
	marshaler_ "github.com/kaydxh/golang/go/runtime/marshaler"
	interceptortcloud3_ "github.com/kaydxh/golang/pkg/middleware/api/tcloud/v3.0"
	interceptortrivialv1_ "github.com/kaydxh/golang/pkg/middleware/api/trivial/v1"
	interceptortrivialv2_ "github.com/kaydxh/golang/pkg/middleware/api/trivial/v2"
	httpinterceptordebug_ "github.com/kaydxh/golang/pkg/middleware/http-middleware/debug"
	httpinterceptorhttp_ "github.com/kaydxh/golang/pkg/middleware/http-middleware/http"
	httpinterceptoropentelemetr_ "github.com/kaydxh/golang/pkg/middleware/http-middleware/opentelemetry"
	httpinterceptorlimiter_ "github.com/kaydxh/golang/pkg/middleware/http-middleware/ratelimiter"
	httpinterceptortimer_ "github.com/kaydxh/golang/pkg/middleware/http-middleware/timer"
)

func WithGatewayMuxOptions(opts ...runtime.ServeMuxOption) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.gatewayMuxOptions = append(c.opts.gatewayMuxOptions, opts...)
	})
}

//now unused, only called by successed response, only append message to response
func WithServerInterceptorsHTTPForwardResponseOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithGatewayMuxOptions(runtime.WithForwardResponseOption(interceptortcloud3_.HTTPForwardResponse)).apply(c)
	})
}

//now unused, only called by successed response
func WithServerInterceptorsTrivialV1HTTPForwardResponseOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithGatewayMuxOptions(runtime.WithForwardResponseOption(interceptortrivialv1_.HTTPForwardResponse)).apply(c)
	})
}

//tcloud api3.0 http response formatter
func WithServerInterceptorsTCloud30HTTPResponseOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithGatewayMuxOptions(
			runtime.WithMarshalerOption(runtime.MIMEWildcard, interceptortcloud3_.NewDefaultJSONPb()),
		).apply(
			c,
		)

		WithGatewayMuxOptions(
			runtime.WithMarshalerOption(binding.MIMEJSON, interceptortcloud3_.NewDefaultJSONPb()),
		).apply(
			c,
		)
	})
}

//trivial api1.0 http response formatter
func WithServerInterceptorsTrivialV1HTTPResponseOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithGatewayMuxOptions(
			runtime.WithMarshalerOption(runtime.MIMEWildcard, marshaler_.NewDefaultJSONPb()),
		).apply(
			c,
		)

		WithGatewayMuxOptions(
			runtime.WithMarshalerOption(binding.MIMEJSON, marshaler_.NewDefaultJSONPb()),
		).apply(
			c,
		)
	})
}

//trivial api2.0 http response formatter
func WithServerInterceptorsTrivialV2HTTPResponseOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithGatewayMuxOptions(
			runtime.WithMarshalerOption(runtime.MIMEWildcard, marshaler_.NewDefaultJSONPb()),
		).apply(
			c,
		)

		WithGatewayMuxOptions(
			runtime.WithMarshalerOption(binding.MIMEJSON, marshaler_.NewDefaultJSONPb()),
		).apply(
			c,
		)
	})
}

// http body proto Marshal
func WithServerInterceptorsHttpBodyProtoOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithGatewayMuxOptions(
			runtime.WithMarshalerOption(
				binding.MIMEPROTOBUF,
				&marshaler_.ProtoMarshaller{},
			),
		).apply(
			c,
		)
	})
}

//HTTP, only called by failed response
func WithServerInterceptorsTCloud30HttpErrorOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithGatewayMuxOptions(runtime.WithErrorHandler(interceptortcloud3_.HTTPError)).apply(c)
	})
}

//HTTP, only called by failed response
func WithServerInterceptorsTrivialV1HttpErrorOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithGatewayMuxOptions(runtime.WithErrorHandler(interceptortrivialv1_.HTTPError)).apply(c)
	})
}

//HTTP, only called by failed response
func WithServerInterceptorsTrivialV2HttpErrorOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithGatewayMuxOptions(runtime.WithErrorHandler(interceptortrivialv2_.HTTPError)).apply(c)
	})
}

// WithHttpPreHandlerInterceptorOptions
func WithHttpPreHandlerInterceptorOptions(
	handlers ...func(w http.ResponseWriter, r *http.Request) error,
) GRPCGatewayOption {

	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.interceptionOptions.httpServerOpts.handlerChain.PreHandlers = append(
			c.opts.interceptionOptions.httpServerOpts.handlerChain.PreHandlers,
			handlers...)
	})
}

// WithHttpHandlerInterceptorOptions
func WithHttpHandlerInterceptorOptions(handlers ...http_.HandlerInterceptor) GRPCGatewayOption {

	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		if handlers != nil {
			c.opts.interceptionOptions.httpServerOpts.handlerChain.Handlers = append(
				c.opts.interceptionOptions.httpServerOpts.handlerChain.Handlers,
				handlers...)
		}
	})
}

// WithHttpPostHandlerInterceptorOptions
func WithHttpPostHandlerInterceptorOptions(
	handlers ...func(w http.ResponseWriter, r *http.Request),
) GRPCGatewayOption {

	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.interceptionOptions.httpServerOpts.handlerChain.PostHandlers = append(
			c.opts.interceptionOptions.httpServerOpts.handlerChain.PostHandlers,
			handlers...)
	})
}

// WithHttpHandlerInterceptorRequestIDAndTraceIDOptions extracts X-Request-ID and X-Traceid
// from HTTP request into context. If absent, X-Request-ID is auto-generated and X-Traceid
// defaults to the same value as X-Request-ID.
func WithHttpHandlerInterceptorRequestIDAndTraceIDOptions() GRPCGatewayOption {
	return WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
		Interceptor: httpinterceptordebug_.RequestIDAndTraceID,
	})
}

// WithHttpHandlerInterceptorRequestIDOptions is an alias for backward compatibility.
// Deprecated: Use WithHttpHandlerInterceptorRequestIDAndTraceIDOptions instead.
func WithHttpHandlerInterceptorRequestIDOptions() GRPCGatewayOption {
	return WithHttpHandlerInterceptorRequestIDAndTraceIDOptions()
}

func WithHttpHandlerInterceptorsTimerOptions() GRPCGatewayOption {
	return WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
		Interceptor: httpinterceptortimer_.ServerInterceptorOfTimer,
	})
}

func WithHttpHandlerInterceptorsMetricOptions() GRPCGatewayOption {
	return WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
		Interceptor: httpinterceptoropentelemetr_.Metric,
	})
}

// WithHttpHandlerInterceptorsTraceOptions adds OpenTelemetry trace interceptor for HTTP requests
func WithHttpHandlerInterceptorsTraceOptions() GRPCGatewayOption {
	return WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
		Interceptor: httpinterceptoropentelemetr_.Trace,
	})
}

func WithHttpHandlerInterceptorInOutputPrinterOptions(enabled bool) GRPCGatewayOption {
	if enabled {
		return WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
			Interceptor: httpinterceptordebug_.InOutputPrinter,
		})
	}

	return WithHttpHandlerInterceptorOptions()
}

func WithHttpHandlerInterceptorInOutputHeaderPrinterOptions(enabled bool) GRPCGatewayOption {
	if enabled {
		return WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
			Interceptor: httpinterceptordebug_.InOutputHeaderPrinter,
		})
	}
	return WithHttpHandlerInterceptorOptions()
}

// timeout
func WithHttpHandlerInterceptorTimeoutOptions(timeout time.Duration) GRPCGatewayOption {
	return WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
		Interceptor: httpinterceptorhttp_.Timeout(timeout),
	})
}

func WithHttpHandlerInterceptorsLimitAllOptions(burst int) GRPCGatewayOption {
	handler := http_.HandlerInterceptor{}
	if burst > 0 {
		handler.Interceptor = httpinterceptorlimiter_.LimitAll(burst).Handler
	}

	return WithHttpHandlerInterceptorOptions(handler)
}

// HTTPQPSLimitConfig HTTP QPS限流配置
type HTTPQPSLimitConfig struct {
	// DefaultQPS 默认QPS限制，0表示不限制
	DefaultQPS float64
	// DefaultBurst 默认突发容量
	DefaultBurst int
	// MaxConcurrency 最大并发数限制，0表示不限制
	// 与QPS限流不同，并发控制限制的是同时处理的请求数，请求完成后令牌会归还
	MaxConcurrency int
	// PathQPS 路径级QPS配置，key为URL路径（如 /api/v1/users）
	PathQPS map[string]float64
	// PathBurst 路径级突发容量配置
	PathBurst map[string]int
	// PathMaxConcurrency 路径级最大并发数配置
	PathMaxConcurrency map[string]int
}

// WithHttpHandlerInterceptorsQPSLimitOptions HTTP QPS限流中间件选项
// 支持全局默认QPS和路径级QPS配置
func WithHttpHandlerInterceptorsQPSLimitOptions(config HTTPQPSLimitConfig) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		// QPS限流
		if config.DefaultQPS > 0 || len(config.PathQPS) > 0 {
			// 创建方法级QPS限流器
			defaultBurst := config.DefaultBurst
			if defaultBurst <= 0 {
				defaultBurst = int(config.DefaultQPS)
			}
			limiter := httpinterceptorlimiter_.NewQPSRateLimiter(config.DefaultQPS, defaultBurst)

			// 设置路径级QPS配置
			for path, qps := range config.PathQPS {
				burst := config.PathBurst[path]
				if burst <= 0 {
					burst = int(qps)
				}
				limiter.SetPathQPS(path, qps, burst)
			}

			WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
				Interceptor: limiter.Handler,
			}).apply(c)
		}

		// 并发控制
		if config.MaxConcurrency > 0 || len(config.PathMaxConcurrency) > 0 {
			limiter := httpinterceptorlimiter_.NewConcurrencyLimiter(config.MaxConcurrency)

			// 设置路径级并发配置
			for path, maxConcurrency := range config.PathMaxConcurrency {
				if maxConcurrency > 0 {
					limiter.SetPathConcurrency(path, maxConcurrency)
				}
			}

			WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
				Interceptor: limiter.Handler,
			}).apply(c)
		}
	})
}

// CleanPath
func WithHttpHandlerInterceptorCleanPathOptions() GRPCGatewayOption {
	return WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
		Interceptor: httpinterceptorhttp_.CleanPath,
	})
}

// recovery
func WithHttpHandlerInterceptorRecoveryOptions() GRPCGatewayOption {
	return WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
		Interceptor: httpinterceptordebug_.Recovery,
	})
}
