/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
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
