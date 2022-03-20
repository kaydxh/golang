package interceptortrace

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func TraceID(w http.ResponseWriter, r *http.Request) {
	requestID := ExtractRequestIdHTTPContext(r)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	SetRequestIdContext(r, requestID)
	fmt.Printf("444444444444444444444444")
}
