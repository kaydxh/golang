package interceptortimer

import (
	"context"

	time_ "github.com/kaydxh/golang/go/time"
	logs_ "github.com/kaydxh/golang/pkg/logs"

	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a new unary server interceptors that timing request
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		tc := time_.New(true)
		logger := logs_.GetLogger(ctx)
		summary := func() {
			tc.Tick(info.FullMethod)
			logger.WithField("method", info.FullMethod).Infof(tc.String())
		}
		defer summary()

		return handler(ctx, req)
	}
}
