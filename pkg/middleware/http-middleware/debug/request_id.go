package interceptordebug

import (
	"net/http"

	"github.com/google/uuid"
	http_ "github.com/kaydxh/golang/go/net/http"
)

func RequestID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := http_.ExtractRequestIdHTTPAndContext(r)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		r = http_.SetRequestIdContext(r, requestID)
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
