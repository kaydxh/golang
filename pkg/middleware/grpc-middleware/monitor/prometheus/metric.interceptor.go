package interceptorprometheus

import (
	"context"
	"fmt"

	grpc_ "github.com/kaydxh/golang/go/net/grpc"
	time_ "github.com/kaydxh/golang/go/time"
	logs_ "github.com/kaydxh/golang/pkg/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptorOfTimer returns a new unary server interceptors that timing request
func UnaryServerInterceptorOfTimer(enabledMetric bool) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		tc := time_.New(true)
		logger := logs_.GetLogger(ctx)
		summary := func() {
			tc.Tick(info.FullMethod)
			if enabledMetric {
				M.durationCost.WithLabelValues(info.FullMethod).Observe(float64(tc.Elapse().Milliseconds()))
			}

			logger.WithField("method", info.FullMethod).Infof(tc.String())
		}
		defer summary()

		return handler(ctx, req)
	}
}

// UnaryServerInterceptorOfCodeMessage returns a new unary server interceptors that timing request
func UnaryServerInterceptorOfCodeMessage(enabledMetric bool) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		peerAddr, _ := grpc_.GetIPFromContext(ctx)
		var (
			resp    interface{}
			code    uint32
			message string
			err     error
		)

		logger := logs_.GetLogger(ctx)
		summary := func() {
			codeMessage := fmt.Sprintf("%d:%s", code, message)
			if enabledMetric {
				metircLabels := map[string]string{
					MetircLabelMethod:      info.FullMethod,
					MetircLabelClientIP:    peerAddr.String(),
					MetircLabelCodeMessage: codeMessage,
				}
				M.calledTotal.With(metircLabels).Inc()
			}

			logger.WithField(
				"method",
				info.FullMethod,
			).WithField(
				"code_messge",
				codeMessage,
			).Infof(
				"called by peer addr: %v",
				peerAddr.String(),
			)
		}
		defer summary()

		resp, err = handler(ctx, req)
		grpcErr, ok := status.FromError(err)
		if ok {
			code = uint32(grpcErr.Code())
			message = grpcErr.Message()

		}
		return resp, err
	}
}