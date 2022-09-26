package interceptorhttp

import (
	"net/http"
	"path"
)

// CleanPath middleware will clean out double slash mistakes from a user's request path.
// For example, if a user requests /users//1 or //users////1 will both be treated as: /users/1
func CleanPath(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var routePath string
		if r.URL.RawPath != "" {
			//RawPath has some like %2F characters.
			routePath = r.URL.RawPath
		} else {
			routePath = r.URL.Path
		}
		r.URL.Path = path.Clean(routePath)

		next.ServeHTTP(w, r)
	})
}
