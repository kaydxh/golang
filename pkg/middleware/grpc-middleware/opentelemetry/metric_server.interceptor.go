package interceptoropentelemetry

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpc_ "github.com/kaydxh/golang/go/net/grpc"
	time_ "github.com/kaydxh/golang/go/time"
	logs_ "github.com/kaydxh/golang/pkg/logs"
	resource_ "github.com/kaydxh/golang/pkg/middleware/resource"
	"google.golang.org/grpc"
)

func UnaryServerMetricInterceptor() grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		return handleMetric(info, resource_.HandlerWithContext[any, any](handler))(ctx, req)
	}
}

func HandleMetric[REQ any, RESP any](
	handler resource_.HandlerWithContext[REQ, RESP],
) resource_.HandlerWithContext[REQ, RESP] {
	return handleMetric(nil, handler)
}

func handleMetric[REQ any, RESP any](
	info *grpc.UnaryServerInfo,
	handler resource_.HandlerWithContext[REQ, RESP],
) resource_.HandlerWithContext[REQ, RESP] {
	return func(ctx context.Context, req REQ) (RESP, error) {

		tc := time_.New(true)

		ctx = resource_.AddKeysContext(ctx, resource_.AddAttrKeysContext, resource_.AddMetricKeysContext)

		resp, err := handler(ctx, req)

		var method string
		if info != nil {
			method = info.FullMethod
		} else {
			method, _ = runtime.RPCMethod(ctx)
		}

		resource_.ReportMetric(ctx,
			resource_.Dimension{
				CalleeMethod: method,
				Error:        err,
			},
			tc.Elapse(),
		)
		tc.Tick(method)

		logger := logs_.GetLogger(ctx)
		peerAddr, _ := grpc_.GetIPFromContext(ctx)
		summary := func() {
			// 每个请求都会触发，开发期 sanity check 用，生产期是噪音。
			// 用 Debug 级别让 level=info 的部署默认屏蔽；OpenTelemetry 已经在
			// 导出 metric，业务侧不需要再依赖访问日志看 QPS/耗时。
			logger.WithField("cost", tc.String()).Debugf(
				"called by peer addr: %v",
				peerAddr.String(),
			)
		}
		defer summary()

		return resp, err
	}
}
