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
package v1

import (
	"context"

	"github.com/google/uuid"
	http_ "github.com/kaydxh/golang/go/net/http"
	reflect_ "github.com/kaydxh/golang/go/reflect"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

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
		ctx = context.WithValue(ctx, http_.DefaultHTTPRequestIDKey, id)
		resp, err := handler(ctx, req)
		// try set requestId to response
		reflect_.TrySetId(req, reflect_.FieldNameSessionId, id)
		// write RequestId to HTTP Header
		if err_ := grpc.SetHeader(ctx, metadata.Pairs(http_.DefaultHTTPRequestIDKey, id)); err_ != nil {
			logrus.WithError(err_).WithField("request_id", id).Warningf("grpc.SendHeader, ignore")
		}
		return resp, err
	}
}
