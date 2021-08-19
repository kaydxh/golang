package tcloud

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func errorToString(err error) string {
	if err == nil {
		return codes.OK.String()
	}

	st, ok := status.FromError(err)
	if !ok {
		return err.Error()
	}

	return st.Code().String()
}

func errorToCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	st, ok := status.FromError(err)
	if !ok {
		return codes.Unknown
	}

	return st.Code()
}

// UnaryServerInterceptorOfError returns a new unary server interceptors
func UnaryServerInterceptorOfError() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		resp, err := handler(ctx, req)
		if err == nil {
			return resp, err
		}

		errResponse := &ErrorResponse{
			Error: &TCloudError{
				Code:    errorToCode(err).String(),
				Message: errorToString(err),
			},
			RequestId: retrieveRequestId(req),
		}
		fmt.Printf("errResponse :%v", errResponse)

		return errResponse, nil
	}
}
