package v1

import (
	"context"

	errors_ "github.com/kaydxh/golang/go/errors"
	reflect_ "github.com/kaydxh/golang/go/reflect"
	"google.golang.org/grpc"
)

// UnaryServerInterceptorOfError returns a new unary server interceptors
func UnaryServerInterceptorOfError() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		resp, err := handler(ctx, req)
		if err == nil {
			return resp, err
		}

		errResponse := &TrivialErrorResponse{
			ErrorCode: int32(errors_.ErrorToCode(err)),
			ErrorMsg:  errors_.ErrorToString(err),
			SessionId: reflect_.RetrieveId(req, reflect_.FieldNameSessionId),
		}

		return errResponse, nil
	}
}
