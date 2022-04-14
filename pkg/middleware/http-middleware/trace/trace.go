package interceptortrace

import (
	"net/http"

	"github.com/google/uuid"
)

func TraceID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := ExtractRequestIdHTTPAndContext(r)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		r = SetRequestIdContext(r, requestID)
		handler.ServeHTTP(w, r)
	})
}

func SetPairsContext(keys []string) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, key := range keys {
				value := ExtractHTTPAndContext(r, key)
				if value != "" {
					r = SetPairContext(r, key, value)
				}
			}

			handler.ServeHTTP(w, r)
		})
	}
}
