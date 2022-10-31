package interceptoropentelemetry

import (
	"context"

	grpc_ "github.com/kaydxh/golang/go/net/grpc"
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

		resp, err = handler(ctx, req)

		attrs := resource_.Attrs(
			resource_.Dimension{
				CalleeMethod: info.FullMethod,
				Error:        err,
			},
		)
		resource_.DefaultMetricMonitor.TotalReqCounter.Add(ctx, 1, attrs...)
		if err != nil {
			resource_.DefaultMetricMonitor.FailCntCounter.Add(ctx, 1, attrs...)
		}

		logger := logs_.GetLogger(ctx)
		peerAddr, _ := grpc_.GetIPFromContext(ctx)
		summary := func() {
			logger.WithField(
				"attrs",
				attrs,
			).Infof(
				"called by peer addr: %v",
				peerAddr.String(),
			)
		}
		defer summary()

		resp, err = handler(ctx, req)
		return resp, err
	}
}
