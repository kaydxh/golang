package ratelimit

import (
	"context"
	"time"

	logs_ "github.com/kaydxh/golang/pkg/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Limiter interface {
	Allow() bool
	AllowFor(timeout time.Duration) bool
	Put()
}

// UnaryServerInterceptor returns a new unary server interceptors that performs request rate limiting.
func UnaryServerInterceptor(limiter Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !limiter.Allow() {
			logger := logs_.GetLogger(ctx)
			err := status.Errorf(
				codes.ResourceExhausted,
				"%s is rejected by grpc_ratelimit middleware, please retry later.",
				info.FullMethod,
			)
			logger.Errorf("%#v", err)
			return nil, err
		}
		defer limiter.Put()
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new stream server interceptor that performs rate limiting on the request.
func StreamServerInterceptor(limiter Limiter) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if !limiter.Allow() {
			logger := logs_.GetLogger(stream.Context())
			err := status.Errorf(
				codes.ResourceExhausted,
				"%s is rejected by grpc_ratelimit middleware, please retry later.",
				info.FullMethod,
			)

			logger.Errorf("%#v", err)
			return err
		}

		defer limiter.Put()
		return handler(srv, stream)
	}
}
