package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	context_ "github.com/kaydxh/golang/go/context"
	errors_ "github.com/kaydxh/golang/go/errors"
	runtime_ "github.com/kaydxh/golang/go/runtime"
	strings_ "github.com/kaydxh/golang/go/strings"

	jsonpb_ "github.com/kaydxh/golang/pkg/protobuf/jsonpb"
	"google.golang.org/protobuf/proto"
)

// HTTPError uses the mux-configured error handler.
func HTTPError(ctx context.Context, mux *runtime.ServeMux,
	marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {

	requestId := context_.ExtractStringFromContext(ctx, RequestIdKey)
	if requestId == "" {
		requestId = strings_.GetStringOrFallback(append(runtime_.GetMetadata(ctx, RequestIdKey), "")...)
	}
	errResponse := &TrivialErrorResponse{
		ErrorCode: int32(errors_.ErrorToCode(err)),
		ErrorMsg:  errors_.ErrorToString(err),
		SessionId: requestId,
	}

	// ForwardResponseMessage forwards the message "resp" from gRPC server to REST client
	runtime.ForwardResponseMessage(ctx, mux, marshaler, w, r, errResponse)
}

//cant not rewrite message, only append message to response
func HTTPForwardResponse(ctx context.Context, r http.ResponseWriter, message proto.Message) error {
	respStruct, err := jsonpb_.MarshaToStructpb(message)
	if err != nil {
		return err
	}

	jb, err := json.Marshal(respStruct)
	if err != nil {
		return fmt.Errorf("jsonpb.Marshal: %v", err)
	}

	r.Write(jb)

	return nil
}
