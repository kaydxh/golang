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
package http

import (
	"context"
	"net/http"
)

const (
	DefaultHTTPRequestIDKey = "X-Request-ID"
)

func ExtractRequestIdHTTPAndContext(r *http.Request) string {
	return ExtractHTTPAndContext(r, DefaultHTTPRequestIDKey)
}

func ExtractHTTPAndContext(r *http.Request, key string) string {
	if value := ExtractFromHTTP(r, key); value != "" {
		return value
	}

	return ExtractFromContext(r.Context(), key)
}

func ExtractFromHTTP(r *http.Request, key string) string {
	if value := r.Header.Get(key); value != "" {
		return value
	}
	if value := r.URL.Query().Get(key); value != "" {
		return value
	}
	if value := r.FormValue(key); value != "" {
		return value
	}
	if value := r.PostFormValue(key); value != "" {
		return value
	}

	return ""
}

func ExtractRequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(DefaultHTTPRequestIDKey).(string); ok {
		return v
	}

	return ""
}

func ExtractFromContext(ctx context.Context, key string) string {
	switch requestIDs := ctx.Value(key).(type) {
	case string:
		if requestIDs != "" {
			return requestIDs
		}
	case []string:
		if len(requestIDs) > 0 {
			return requestIDs[0]
		}
	}

	return ""

}

func SetPairContext(r *http.Request, key, value string) *http.Request {
	ctx := context.WithValue(r.Context(), key, value)
	r = r.WithContext(ctx)

	return r
}

func SetRequestIdContext(r *http.Request, requestID string) *http.Request {
	return SetPairContext(r, DefaultHTTPRequestIDKey, requestID)
}
