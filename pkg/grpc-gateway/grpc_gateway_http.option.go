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
	httpinterceptoropentelemetr_ "github.com/kaydxh/golang/pkg/middleware/http-middleware/monitor/opentelemetry"
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

// WithHttpHandlerInterceptorRequestIDOptions
func WithHttpHandlerInterceptorRequestIDOptions() GRPCGatewayOption {
	return WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
		Interceptor: httpinterceptordebug_.RequestID,
	})
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
