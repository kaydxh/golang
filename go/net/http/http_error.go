package http

import (
	"fmt"
	"net/http"
)

func ErrorFromHttp(code int) error {

	if code >= http.StatusOK && code < http.StatusBadRequest {
		return nil
	}

	return fmt.Errorf("unexpected http status code: %v", code)
}
