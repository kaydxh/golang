package interceptoropentelemetry

import (
	"context"

	grpc_ "github.com/kaydxh/golang/go/net/grpc"
	time_ "github.com/kaydxh/golang/go/time"
	logs_ "github.com/kaydxh/golang/pkg/logs"
	resource_ "github.com/kaydxh/golang/pkg/middleware/resource"
	"google.golang.org/grpc"
)

func UnaryServerMetricInterceptor() grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		var (
			resp interface{}
			err  error
		)

		tc := time_.New(true)
		resp, err = handler(ctx, req)

		resource_.ReportMetric(ctx,
			resource_.Dimension{
				CalleeMethod: info.FullMethod,
				Error:        err,
			},
			tc.Elapse(),
		)
		tc.Tick(info.FullMethod)

		logger := logs_.GetLogger(ctx)
		peerAddr, _ := grpc_.GetIPFromContext(ctx)
		summary := func() {
			logger.WithField("cost", tc.String()).Infof(
				"called by peer addr: %v",
				peerAddr.String(),
			)
		}
		defer summary()

		return resp, err
	}
}
