package grpcgateway_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	gw_ "github.com/kaydxh/golang/pkg/grpc-gateway"

	"github.com/kaydxh/golang/pkg/grpc-gateway/date"
	//"github.com/kaydxh/sea/api/openapi-spec/v1/date"
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
