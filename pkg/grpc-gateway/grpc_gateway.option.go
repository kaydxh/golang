package grpcgateway

import (
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
