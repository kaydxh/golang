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
	httpinterceptordebug_ "github.com/kaydxh/golang/pkg/middleware/http-middleware/debug"
	httpinterceptorhttp_ "github.com/kaydxh/golang/pkg/middleware/http-middleware/http"
	httpinterceptorprometheus_ "github.com/kaydxh/golang/pkg/middleware/http-middleware/monitor/prometheus"
	httpinterceptorlimiter_ "github.com/kaydxh/golang/pkg/middleware/http-middleware/ratelimiter"
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
			runtime.WithMarshalerOption(runtime.MIMEWildcard, interceptortrivialv1_.NewDefaultJSONPb()),
		).apply(
			c,
		)

		WithGatewayMuxOptions(
			runtime.WithMarshalerOption(binding.MIMEJSON, interceptortrivialv1_.NewDefaultJSONPb()),
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

func WithHttpHandlerInterceptorsTimerOptions(enabledMetric bool) GRPCGatewayOption {
	return WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
		Interceptor: httpinterceptorprometheus_.InterceptorOfTimer(enabledMetric),
	})
}

func WithHttpHandlerInterceptorInOutputPrinterOptions() GRPCGatewayOption {
	return WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
		Interceptor: httpinterceptordebug_.InOutputPrinter,
	})
}

func WithHttpHandlerInterceptorInOutputHeaderPrinterOptions() GRPCGatewayOption {
	return WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
		Interceptor: httpinterceptordebug_.InOutputHeaderPrinter,
	})
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
