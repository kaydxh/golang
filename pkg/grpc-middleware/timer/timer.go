package interceptortimer

import (
	"context"

	time_ "github.com/kaydxh/golang/go/time"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a new unary server interceptors that timing request
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		tc := time_.New(true)
		summary := func() {
			tc.Tick(info.FullMethod)
			logrus.Infof(tc.String())
		}
		defer summary()
		return handler(ctx, req)
	}
}
