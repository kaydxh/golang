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
