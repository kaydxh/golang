package interceptorhttp

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	logs_ "github.com/kaydxh/golang/pkg/logs"
)

// DefaultTimeoutHandler is a convenient timeout handler which
// simply returns "504 Service timeout".
var DefaultTimeoutHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGatewayTimeout)
		w.Write([]byte("Service timeout"))
	})

type timeoutWriter struct {
	w http.ResponseWriter

	status int
	buf    *bytes.Buffer
}

func (tw timeoutWriter) Header() http.Header {
	return tw.w.Header()
}

func (tw *timeoutWriter) WriteHeader(status int) {
	tw.status = status
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	if tw.status == 0 {
		tw.status = http.StatusOK
	}

	return tw.buf.Write(b)
}

func (tw *timeoutWriter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

// Timeout is a middleware that cancels ctx after a given timeout and return
// a 504 Gateway Timeout error to the client.
func Timeout(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			logger := logs_.GetLogger(r.Context())
			background := r.Context()
			tCtx, tCancel := context.WithTimeout(background, timeout)
			cCtx, cCancel := context.WithCancel(background)
			r = r.WithContext(cCtx)
			defer tCancel()

			tw := &timeoutWriter{w: w, buf: bytes.NewBuffer(nil)}
			go func() {
				//next.ServeHTTP(tw, r)
				next.ServeHTTP(w, r)
				cCancel()
			}()

			select {
			case <-cCtx.Done():
				fmt.Printf("----normal\n")
			case <-tCtx.Done():
				if err := tCtx.Err(); err == context.DeadlineExceeded {
					logger.WithField("method", fmt.Sprintf("%s %s", r.Method, r.URL.Path)).Errorf("server %v timeout", timeout)
					cCancel()
					//		DefaultTimeoutHandler.ServeHTTP(w, r)
					tw.ServeHTTP(w, r)
				}
			}
		}
		return http.HandlerFunc(fn)
	}
}
