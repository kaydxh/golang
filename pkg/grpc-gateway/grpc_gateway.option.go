package grpcgateway

import (
	"github.com/gin-gonic/gin/binding"
	interceptorlogrus_ "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	interceptortcloud_ "github.com/kaydxh/golang/pkg/grpc-middleware/api/tcloud/v3.0"
	interceptortimer_ "github.com/kaydxh/golang/pkg/grpc-middleware/timer"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

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

func WithGatewayMuxOptions(opts ...runtime.ServeMuxOption) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.gatewayMuxOptions = append(c.opts.gatewayMuxOptions, opts...)
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

func WithServerUnaryInterceptorsTimerOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(interceptortimer_.UnaryServerInterceptor()).apply(c)
	})
}

func WithServerUnaryInterceptorsRequestIdOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(interceptortcloud_.UnaryServerInterceptorOfRequestId()).apply(c)
	})
}

func WithServerUnaryInterceptorsErrorOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(interceptortcloud_.UnaryServerInterceptorOfError()).apply(c)
		//WithServerStreamInterceptorsOptions(interceptortcloud_.StreamServerInterceptor()).apply(c)
	})
}

//HTTP, only called by failed response
func WithServerInterceptorsHttpErrorOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithGatewayMuxOptions(runtime.WithErrorHandler(interceptortcloud_.HTTPError)).apply(c)
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
