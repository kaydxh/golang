package middleware

import (
	"context"

	interceptordebug_ "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/debug"
	interceptoropentelemetry_ "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/monitor/opentelemetry"
)

// the function need to called by controller
func LocalMiddlewareWrap[REQ any, RESP any](handler func(ctx context.Context, req REQ) (RESP, error)) func(ctx context.Context, req REQ) (RESP, error) {
	return interceptordebug_.HandleReuestId(
		interceptordebug_.HandleInOutputPrinter(
			interceptoropentelemetry_.HandleMetric(
				handler,
			)))
}