package interceptortrace

import (
	"net/http"

	"github.com/google/uuid"
)

func TraceID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := ExtractRequestIdHTTPContext(r)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		r = SetRequestIdContext(r, requestID)
		handler.ServeHTTP(w, r)
	})
}
