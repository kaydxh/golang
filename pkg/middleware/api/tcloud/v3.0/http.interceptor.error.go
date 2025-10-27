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
package tcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	context_ "github.com/kaydxh/golang/go/context"
	errors_ "github.com/kaydxh/golang/go/errors"
	http_ "github.com/kaydxh/golang/go/net/http"
	runtime_ "github.com/kaydxh/golang/go/runtime"
	strings_ "github.com/kaydxh/golang/go/strings"
	jsonpb_ "github.com/kaydxh/golang/pkg/protobuf/jsonpb"
	"github.com/sirupsen/logrus"

	"google.golang.org/protobuf/proto"
)

// HTTPError uses the mux-configured error handler.
func HTTPError(ctx context.Context, mux *runtime.ServeMux,
	marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {

	requestId := context_.ExtractStringFromContext(ctx, http_.DefaultHTTPRequestIDKey)
	if requestId == "" {
		requestId = strings_.GetStringOrFallback(
			append(runtime_.GetMetadata(ctx, http_.DefaultHTTPRequestIDKey), "")...)
	}

	errResponse := &ErrorResponse{
		Error: &TCloudError{
			Code:    errors_.ErrorToCodeString(err),
			Message: errors_.ErrorToString(err),
		},
		RequestId: requestId,
	}
	logrus.WithField("request_id", requestId).WithField("response", errResponse).Errorf("send")

	// ForwardResponseMessage forwards the message "resp" from gRPC server to REST client
	runtime.ForwardResponseMessage(ctx, mux, marshaler, w, r, errResponse)
}

//cant not rewrite message, only append message to response
func HTTPForwardResponse(ctx context.Context, r http.ResponseWriter, message proto.Message) error {
	respStruct, err := jsonpb_.MarshaToStructpb(message)
	if err != nil {
		return err
	}
	errResponse := &TCloudResponse{
		Response: respStruct,
	}

	jb, err := json.Marshal(errResponse)
	if err != nil {
		return fmt.Errorf("jsonpb.Marshal: %v", err)
	}

	r.Write(jb)

	return nil
}
