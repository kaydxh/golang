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
