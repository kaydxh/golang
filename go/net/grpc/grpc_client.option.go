package grpc

import (
	"time"
)

func WithMaxMsgSize(maxMsgSize int) GrpcClientOption {
	return GrpcClientOptionFunc(func(c *GrpcClient) {
		c.opts.maxMsgSize = maxMsgSize
	})
}

func WithConnectionTimeout(connectionTimeout time.Duration) GrpcClientOption {
	return GrpcClientOptionFunc(func(c *GrpcClient) {
		c.opts.connectionTimeout = connectionTimeout
	})
}
