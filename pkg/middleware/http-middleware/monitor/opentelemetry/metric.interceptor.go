package interceptoropentelemetry

import (
	"fmt"
	"net/http"

	http_ "github.com/kaydxh/golang/go/net/http"
	logs_ "github.com/kaydxh/golang/pkg/logs"
	resource_ "github.com/kaydxh/golang/pkg/middleware/resource"
)

func Metric(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)

		ctx := r.Context()
		func() {
			attrs := resource_.Attrs(
				resource_.Dimension{
					CalleeMethod: fmt.Sprintf("%v %v", r.Method, r.URL.Path),
					Error:        nil,
				},
			)
			resource_.DefaultMetricMonitor.TotalReqCounter.Add(ctx, 1, attrs...)
		}()

		logger := logs_.GetLogger(ctx)
		peerAddr, _ := http_.GetIPFromRequest(r)
		summary := func() {
			logger.Infof(
				"called by peer addr: %v",
				peerAddr.String(),
			)
		}
		defer summary()
	})

}
