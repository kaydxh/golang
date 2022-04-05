package grpcgateway

import (
	"net/http"

	"github.com/gin-gonic/gin/binding"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	http_ "github.com/kaydxh/golang/go/net/http"
	interceptortcloud_ "github.com/kaydxh/golang/pkg/grpc-middleware/api/tcloud/v3.0"
	httpinterceptorprometheus_ "github.com/kaydxh/golang/pkg/middleware/http-middleware/monitor/prometheus"
	httpinterceptortrace_ "github.com/kaydxh/golang/pkg/middleware/http-middleware/trace"
)

func WithGatewayMuxOptions(opts ...runtime.ServeMuxOption) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.gatewayMuxOptions = append(c.opts.gatewayMuxOptions, opts...)
	})
}

//now unused, only called by successed response
func WithServerInterceptorsHTTPForwardResponseOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithGatewayMuxOptions(runtime.WithForwardResponseOption(interceptortcloud_.HTTPForwardResponse)).apply(c)
	})
}

//tcloud api3.0 http response formatter
func WithServerInterceptorsTCloud30HTTPResponseOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithGatewayMuxOptions(
			runtime.WithMarshalerOption(runtime.MIMEWildcard, interceptortcloud_.NewDefaultJSONPb()),
		).apply(
			c,
		)

		WithGatewayMuxOptions(
			runtime.WithMarshalerOption(binding.MIMEJSON, interceptortcloud_.NewDefaultJSONPb()),
		).apply(
			c,
		)
	})
}

//HTTP, only called by failed response
func WithServerInterceptorsHttpErrorOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithGatewayMuxOptions(runtime.WithErrorHandler(interceptortcloud_.HTTPError)).apply(c)
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
		c.opts.interceptionOptions.httpServerOpts.handlerChain.Handlers = append(
			c.opts.interceptionOptions.httpServerOpts.handlerChain.Handlers,
			handlers...)
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

// WithHttpHandlerInterceptorTraceIDOptions
func WithHttpHandlerInterceptorTraceIDOptions() GRPCGatewayOption {
	return WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
		Interceptor: httpinterceptortrace_.TraceID,
	})
}

func WithHttpHandlerInterceptorsTimerOptions(enabledMetric bool) GRPCGatewayOption {
	return WithHttpHandlerInterceptorOptions(http_.HandlerInterceptor{
		Interceptor: httpinterceptorprometheus_.InterceptorOfTimer(enabledMetric),
	})
}
