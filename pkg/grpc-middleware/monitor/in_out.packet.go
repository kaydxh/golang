package interceptorsmonitor

import (
	"context"
	"fmt"

	reflect_ "github.com/kaydxh/golang/go/reflect"
	"github.com/sirupsen/logrus"
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

// UnaryServerInterceptorOfInOutPacket log in-output packet
func UnaryServerInterceptorOfInOutPacket() grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		if req != nil {
			//logrus.WithField("request", &JsonpbMarshaller{req.(proto.Message)}).Info("recv")
			logrus.WithField("request", reflect_.TruncateBytes(proto.Clone(req.(proto.Message)))).Info("recv")
		}

		resp, err := handler(ctx, req)
		if resp != nil {
			//logrus.WithField("response", &JsonpbMarshaller{resp.(proto.Message)}).Info("send")
			logrus.WithField("response", reflect_.TruncateBytes(proto.Clone(resp.(proto.Message)))).Info("send")
		}

		return resp, err
	}
}
