package interceptoropentelemetry

import (
	"fmt"
	"net/http"

	errors_ "github.com/kaydxh/golang/go/errors"
	http_ "github.com/kaydxh/golang/go/net/http"
	time_ "github.com/kaydxh/golang/go/time"
	logs_ "github.com/kaydxh/golang/pkg/logs"
	resource_ "github.com/kaydxh/golang/pkg/middleware/resource"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // default status code
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Status() int {
	return rw.statusCode
}

func Metric(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tc := time_.New(true)

		// Wrap ResponseWriter to capture status code
		wrapped := newResponseWriter(w)
		next.ServeHTTP(wrapped, r)

		ctx := r.Context()
		calleeMethod := fmt.Sprintf("%v %v", r.Method, r.URL.Path)

		// Convert HTTP status code to error for metric reporting
		// Only 2xx status codes are considered successful
		var err error
		statusCode := wrapped.Status()
		if statusCode < 200 || statusCode >= 300 {
			err = errors_.Errorf(statusCode, http.StatusText(statusCode))
		}

		resource_.ReportMetric(ctx,
			resource_.Dimension{
				CalleeMethod: calleeMethod,
				Error:        err,
			},
			tc.Elapse(),
		)
		tc.Tick(calleeMethod)

		logger := logs_.GetLogger(ctx)
		peerAddr, _ := http_.GetIPFromRequest(r)
		summary := func() {
			// HEAD/OPTIONS 是 health check / CORS 预检，频次高、无业务意义，
			// 降到 Debug 避免在 level=info 部署下刷屏；业务请求仍打 Info 便于
			// 看请求 QPS/耗时（OpenTelemetry 已经在导出 metric，但 access log
			// 仍是开发期最直观的 sanity check）。
			entry := logger.WithField("cost", tc.String()).WithField("status", statusCode)
			if r.Method == http.MethodHead || r.Method == http.MethodOptions {
				entry.Debugf("called by peer addr: %v", peerAddr.String())
			} else {
				entry.Infof("called by peer addr: %v", peerAddr.String())
			}
		}
		defer summary()
	})
}
