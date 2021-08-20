package tcloud

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	strings_ "github.com/kaydxh/golang/go/strings"
)

// HTTPError uses the mux-configured error handler.
func HTTPError(ctx context.Context, mux *runtime.ServeMux,
	marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	errResponse := &ErrorResponse{
		Error: &TCloudError{
			Code:    errorToCode(err).String(),
			Message: errorToString(err),
		},
		RequestId: strings_.GetStringOrFallback(GetMetadata(ctx, RequestIdKey), ""),
	}
	// ForwardResponseMessage forwards the message "resp" from gRPC server to REST client
	runtime.ForwardResponseMessage(ctx, mux, marshaler, w, r, errResponse)
}
