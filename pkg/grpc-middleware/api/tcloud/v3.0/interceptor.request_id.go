package tcloud

import (
	"context"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RequestIdKey is metadata key name for request ID
var RequestIdKey = "X-Request-ID"

func GetMetadata(ctx context.Context, key string) []string {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok || md.HeaderMD == nil {
		return nil
	}

	return md.HeaderMD.Get(key)
}

// UnaryServerInterceptorOfRequestId returns a new unary server interceptors with tags in context with request_id.
func UnaryServerInterceptorOfRequestId() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// retrieve requestId from request
		id := retrieveRequestId(req)
		if id == "" {
			//if id is empty, set new uuid to request
			id = uuid.NewString()
			trySetRequestId(req, id)
		}

		//set "X-Request-ID" to context
		ctx = context.WithValue(ctx, RequestIdKey, id)
		resp, err := handler(ctx, req)
		// try set requestId to response
		trySetRequestId(resp, id)
		// write RequestId to HTTP Header
		if err_ := grpc.SetHeader(ctx, metadata.Pairs(RequestIdKey, id)); err_ != nil {
			logrus.WithError(err_).WithField("request_id", id).Warningf("grpc.SendHeader, ignore")
		}
		return resp, err
	}
}
