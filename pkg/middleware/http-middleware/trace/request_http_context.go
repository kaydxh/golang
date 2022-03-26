package interceptortrace

import (
	"context"
	"net/http"
)

const (
	DefaultHTTPRequestIDKey = "X-Request-ID"
)

func ExtractRequestIdHTTPContext(r *http.Request) string {
	return ExtractHTTPContext(r, DefaultHTTPRequestIDKey)
}

func ExtractHTTPContext(r *http.Request, key string) string {
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

func SetRequestIdContext(r *http.Request, requestID string) *http.Request {
	ctx := context.WithValue(r.Context(), DefaultHTTPRequestIDKey, requestID)
	r = r.WithContext(ctx)

	return r
}
