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
package grpcgateway_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	gw_ "github.com/kaydxh/golang/pkg/grpc-gateway"

	"github.com/kaydxh/golang/pkg/grpc-gateway/date"
	"google.golang.org/grpc"
)

type Controller struct {

	// Embed the unimplemented server
	date.UnimplementedDateServiceServer
}

// 日期查询
func (c *Controller) Now(
	_ context.Context,
	req *date.DateRequest,
) (resp *date.DateResponse, err error) {
	fmt.Println(">>>>>NOW")
	return &date.DateResponse{
		RequestId: req.GetRequestId(),
		Date:      time.Now().String(),
	}, nil
}

func TestGRPCGateway(t *testing.T) {
	s := Controller{}
	gwServer := gw_.NewGRPCGateWay(":10002")
	gwServer.RegisterGRPCHandler(func(srv *grpc.Server) {
		date.RegisterDateServiceServer(srv, &s)
	})

	err := gwServer.RegisterHTTPHandler(
		context.Background(),
		func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
			return date.RegisterDateServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
		},
	)

	if err != nil {
		t.Fatalf("failed to RegisterHTTPHandler got: %s", err)
	}

	gwServer.ListenAndServe()
}
