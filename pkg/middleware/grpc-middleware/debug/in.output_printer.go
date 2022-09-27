/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
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
