package interceptordebug

import (
	"context"
	"fmt"

	reflect_ "github.com/kaydxh/golang/go/reflect"
	logs_ "github.com/kaydxh/golang/pkg/logs"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type JsonpbMarshaller struct {
	proto.Message
}

func (j *JsonpbMarshaller) MarshalJson() ([]byte, error) {
	data, err := protojson.Marshal(j.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshaJson: %v", err)
	}

	return data, nil
}

//  UnaryServerInterceptorOfInOutputPrinter log in-output packet
func UnaryServerInterceptorOfInOutputPrinter() grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		logger := logs_.GetLogger(ctx)
		if req != nil {
			//logrus.WithField("request", &JsonpbMarshaller{req.(proto.Message)}).Info("recv")
			logger.WithField("request", reflect_.TruncateBytes(proto.Clone(req.(proto.Message)))).Info("recv")
		}

		resp, err := handler(ctx, req)
		if resp != nil {
			//logrus.WithField("response", &JsonpbMarshaller{resp.(proto.Message)}).Info("send")
			logger.WithField("response", reflect_.TruncateBytes(proto.Clone(resp.(proto.Message)))).Info("send")
		}

		return resp, err
	}
}
