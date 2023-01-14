/*
 *Copyright (c) 2023, kaydxh
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

	reflect_ "github.com/kaydxh/golang/go/reflect"
	strings_ "github.com/kaydxh/golang/go/strings"
	logs_ "github.com/kaydxh/golang/pkg/logs"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func UnaryClientInterceptorOfInOutputPrinter(filterMethods ...string) grpc.UnaryClientInterceptor {

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		logger := logs_.GetLogger(ctx)

		enablePrint := !strings_.SliceContainsCaseInSensitive(filterMethods, method)
		if enablePrint {
			reqProto, ok := any(req).(proto.Message)
			if ok {
				logger.WithField("request", reflect_.TruncateBytes(proto.Clone(reqProto))).Info("send")
			}
		}

		err := invoker(ctx, method, req, reply, cc, opts...)

		if enablePrint {
			respProto, ok := any(reply).(proto.Message)
			if ok {
				logger.WithField("response", reflect_.TruncateBytes(proto.Clone(respProto))).Info("recv")
			}
		}

		return err
	}
}
