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
package interceptordebug

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	http_ "github.com/kaydxh/golang/go/net/http"
	logs_ "github.com/kaydxh/golang/pkg/logs"
)

// isProbeRequest 判断是否为 health check / CORS 预检之类的探活请求。
// 这类请求频次高、无业务体，访问日志降到 Debug 避免刷屏；业务请求保持 Info。
func isProbeRequest(method string) bool {
	return method == http.MethodHead || method == http.MethodOptions
}

func InOutputPrinter(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logs_.GetLogger(r.Context())
		ww := http_.NewResponseWriterWrapper(w)

		calleeMethod := fmt.Sprintf("%v %v", r.Method, r.URL.Path)
		probe := isProbeRequest(r.Method)

		defer func() {
			entry := logger.WithField("method", calleeMethod).WithField("response", ww.String())
			if probe {
				entry.Debug("send")
			} else {
				entry.Info("send")
			}
		}()
		if r != nil {
			buf, err := io.ReadAll(r.Body)
			if err != nil {
				return
			}
			rdr := io.NopCloser(bytes.NewBuffer(buf))
			r.Body = rdr
			entry := logger.WithField("method", calleeMethod).WithField("request", string(buf))
			if probe {
				entry.Debug("recv")
			} else {
				entry.Info("recv")
			}

		}

		handler.ServeHTTP(ww, r)

	})
}

func InOutputHeaderPrinter(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logs_.GetLogger(r.Context())
		probe := isProbeRequest(r.Method)

		recvEntry := logger.WithField("request headers", r.Header)
		if probe {
			recvEntry.Debug("recv")
		} else {
			recvEntry.Info("recv")
		}

		defer func() {
			sendEntry := logger.WithField("response headers", w.Header())
			if probe {
				sendEntry.Debug("send")
			} else {
				sendEntry.Info("send")
			}
		}()

		handler.ServeHTTP(w, r)
	})
}
