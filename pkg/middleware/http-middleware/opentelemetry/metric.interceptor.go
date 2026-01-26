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
			logger.WithField("cost", tc.String()).WithField("status", statusCode).Infof(
				"called by peer addr: %v",
				peerAddr.String(),
			)
		}
		defer summary()
	})
}
