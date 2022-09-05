package ratelimit

import (
	"fmt"
	"net/http"
	"time"

	rate_ "github.com/kaydxh/golang/go/time/rate"
	logs_ "github.com/kaydxh/golang/pkg/logs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Limiter interface {
	Allow() bool
	AllowFor(timeout time.Duration) bool
	Put()
}

type RateLimiter struct {
	Limiter
}

/*
func LimitAll(burst int) http.Handler {
	limiter := rate_.NewLimiter(burst)
	rl := &RateLimiter{
		Limiter: limiter,
	}
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	return rl.Handler(h)
}
*/

func LimitAll(burst int) *RateLimiter {
	limiter := rate_.NewLimiter(burst)
	rl := &RateLimiter{
		Limiter: limiter,
	}
	//h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	return rl
}

func (l *RateLimiter) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !l.Allow() {
			logger := logs_.GetLogger(r.Context())
			err := status.Errorf(
				codes.ResourceExhausted,
				"%s is rejected by http_ratelimit middleware, please retry later.",
				fmt.Sprintf("%s %s", r.Method, r.URL.Path),
			)
			logger.Errorf("%v", err.Error())
			return
		}
		defer l.Put()

		handler.ServeHTTP(w, r)
	})
}
