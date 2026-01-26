package interceptoropentelemetry

import (
	"fmt"
	"net/http"

	http_ "github.com/kaydxh/golang/go/net/http"
	time_ "github.com/kaydxh/golang/go/time"
	logs_ "github.com/kaydxh/golang/pkg/logs"
	resource_ "github.com/kaydxh/golang/pkg/middleware/resource"
)

func Metric(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tc := time_.New(true)
		next.ServeHTTP(w, r)

		ctx := r.Context()
		calleeMethod := fmt.Sprintf("%v %v", r.Method, r.URL.Path)
		resource_.ReportMetric(ctx,
			resource_.Dimension{
				CalleeMethod: calleeMethod,
				Error:        nil,
			},
			tc.Elapse(),
		)
		tc.Tick(calleeMethod)

		logger := logs_.GetLogger(ctx)
		peerAddr, _ := http_.GetIPFromRequest(r)
		summary := func() {
			logger.WithField("cost", tc.String()).Infof(
				"called by peer addr: %v",
				peerAddr.String(),
			)
		}
		defer summary()
	})

}
