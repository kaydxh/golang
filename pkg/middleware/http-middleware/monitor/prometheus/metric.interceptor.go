package interceptorprometheus

import (
	"net/http"

	time_ "github.com/kaydxh/golang/go/time"
	logs_ "github.com/kaydxh/golang/pkg/logs"
	prometheus_ "github.com/kaydxh/golang/pkg/monitor/prometheus"
)

func InterceptorOfTimer(enabledMetric bool) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tc := time_.New(true)
			logger := logs_.GetLogger(r.Context())
			summary := func() {
				tc.Tick(r.Method)
				if enabledMetric {
					prometheus_.M.DurationCost.WithLabelValues(
						r.Method,
					).Observe(
						float64(tc.Elapse().Milliseconds()),
					)
				}

				logger.WithField("method", r.Method).Infof(tc.String())
			}
			defer summary()

			handler.ServeHTTP(w, r)
		})

	}

}
