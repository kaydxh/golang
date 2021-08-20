package tcloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// HTTPError uses the mux-configured error handler.
func HTTPError(ctx context.Context, mux *runtime.ServeMux,
	marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	fmt.Printf("err: %v", err)
	errResponse := &ErrorResponse{
		Error: &TCloudError{
			Code:    errorToCode(err).String(),
			Message: errorToString(err),
		},
		RequestId: retrieveRequestId(r),
	}
	// ForwardResponseMessage forwards the message "resp" from gRPC server to REST client
	runtime.ForwardResponseMessage(ctx, mux, marshaler, w, r, errResponse)
}
