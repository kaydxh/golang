package grpcgateway

import (
	interceptorlogrus_ "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func WithServerOptions(opts []grpc.ServerOption) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.serverOptions = opts
	})
}

func WithClientDialOptions(opts []grpc.DialOption) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.clientDialOptions = opts
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

func WithServerUnaryInterceptorsLogrusOptions(
	logger *logrus.Logger,
) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		l := logrus.NewEntry(logger)
		WithServerUnaryInterceptorsOptions(interceptorlogrus_.UnaryServerInterceptor(l))
		WithServerStreamInterceptorsOptions(interceptorlogrus_.StreamServerInterceptor(l))
	})
}
