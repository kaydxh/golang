package interceptortrace

import (
	"context"
	"net/http"
)

const (
	DefaulXHTTPRequestIDKey = "X-Request-ID"
)

func ExtractRequestIdHTTPContext(r *http.Request) string {
	return ExtractHTTPContext(r, DefaulXHTTPRequestIDKey)
}

func ExtractHTTPContext(r *http.Request, key string) string {
	if requestID := r.Header.Get(key); requestID != "" {
		return requestID
	}
	if requestID := r.URL.Query().Get(key); requestID != "" {
		return requestID
	}
	if requestID := r.FormValue(key); requestID != "" {
		return requestID
	}
	if requestID := r.PostFormValue(key); requestID != "" {
		return requestID
	}

	return ExtractRequestIdContext(r.Context())
}

func ExtractRequestIdContext(ctx context.Context) string {
	switch requestIDs := ctx.Value(DefaulXHTTPRequestIDKey).(type) {
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

func SetRequestIdContext(r *http.Request, requestID string) {
	ctx := context.WithValue(r.Context(), DefaulXHTTPRequestIDKey, requestID)
	r = r.WithContext(ctx)

	return
}
