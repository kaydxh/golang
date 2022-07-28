package v1

import (
	"context"

	"github.com/google/uuid"
	reflect_ "github.com/kaydxh/golang/go/reflect"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RequestIdKey is metadata key name for request ID
var RequestIdKey = "X-Request-ID"

// UnaryServerInterceptorOfRequestId returns a new unary server interceptors with tags in context with request_id.
func UnaryServerInterceptorOfRequestId() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// retrieve requestId from request
		id := reflect_.RetrieveId(req, reflect_.FieldNameSessionId)
		if id == "" {
			//if id is empty, set new uuid to request
			id = uuid.NewString()
			reflect_.TrySetId(req, reflect_.FieldNameSessionId, id)
		}

		//set "X-Request-ID" to context
		ctx = context.WithValue(ctx, RequestIdKey, id)
		resp, err := handler(ctx, req)
		// try set requestId to response
		reflect_.TrySetId(req, reflect_.FieldNameSessionId, id)
		// write RequestId to HTTP Header
		if err_ := grpc.SetHeader(ctx, metadata.Pairs(RequestIdKey, id)); err_ != nil {
			logrus.WithError(err_).WithField("request_id", id).Warningf("grpc.SendHeader, ignore")
		}
		return resp, err
	}
}
