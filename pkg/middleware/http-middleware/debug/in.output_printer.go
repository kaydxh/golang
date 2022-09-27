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
	"io/ioutil"
	"net/http"

	http_ "github.com/kaydxh/golang/go/net/http"
	logs_ "github.com/kaydxh/golang/pkg/logs"
)

func InOutputPrinter(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logs_.GetLogger(r.Context())
		ww := http_.NewResponseWriterWrapper(w)
		defer func() {
			logger.WithField("response", ww.String()).Info("send")
		}()
		if r != nil {
			buf, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return
			}
			rdr := ioutil.NopCloser(bytes.NewBuffer(buf))
			r.Body = rdr
			logger.WithField("request", string(buf)).Info("recv")

		}

		handler.ServeHTTP(ww, r)

	})
}

func InOutputHeaderPrinter(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logs_.GetLogger(r.Context())
		logger.WithField("request headers", r.Header).Info("recv")

		defer func() {
			logger.WithField("response headers", w.Header()).Info("send")
		}()

		handler.ServeHTTP(w, r)
	})
}
