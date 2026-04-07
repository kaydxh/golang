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
	"net/http"

	"github.com/google/uuid"
	http_ "github.com/kaydxh/golang/go/net/http"
)

// RequestID extracts X-Request-ID from HTTP request, generates one if absent, and sets it into context.
// Deprecated: Use RequestIDAndTraceID instead.
func RequestID(handler http.Handler) http.Handler {
	return RequestIDAndTraceID(handler)
}

// RequestIDAndTraceID extracts X-Request-ID and X-Traceid from HTTP request into context.
// If X-Request-ID is absent, a new UUID is generated.
// If X-Traceid is absent, it defaults to the same value as X-Request-ID.
func RequestIDAndTraceID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := http_.ExtractRequestIdHTTPAndContext(r)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		r = http_.SetRequestIdContext(r, requestID)

		traceID := http_.ExtractTraceIdHTTPAndContext(r)
		if traceID == "" {
			traceID = requestID
		}
		r = http_.SetTraceIdContext(r, traceID)

		handler.ServeHTTP(w, r)
	})
}

func SetPairsContext(keys []string) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, key := range keys {
				value := http_.ExtractHTTPAndContext(r, key)
				if value != "" {
					r = http_.SetPairContext(r, key, value)
				}
			}

			handler.ServeHTTP(w, r)
		})
	}
}
