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

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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
		return HandleInOutputPrinter(info, handler)(ctx, req)
	}
}

func HandleInOutputPrinter[REQ any, RESP any](info *grpc.UnaryServerInfo, handler func(ctx context.Context, req REQ) (RESP, error)) func(ctx context.Context, req REQ) (RESP, error) {
	return func(ctx context.Context, req REQ) (RESP, error) {
		logger := logs_.GetLogger(ctx)

		var method string
		if info != nil {
			method = info.FullMethod
		} else {
			method, _ = runtime.RPCMethod(ctx)
		}

		reqProto, ok := any(req).(proto.Message)
		if ok {
			logger.WithField("method", method).WithField("request", reflect_.TruncateBytes(proto.Clone(reqProto))).Info("recv")
		}

		resp, err := handler(ctx, req)
		respProto, ok := any(resp).(proto.Message)
		if ok {
			logger.WithField("method", method).WithField("response", reflect_.TruncateBytes(proto.Clone(respProto))).Info("send")
		}

		return resp, err
	}
}
