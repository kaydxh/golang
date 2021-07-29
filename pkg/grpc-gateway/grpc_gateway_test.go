package grpcgateway_test

import (
	"context"
	"fmt"
	"testing"
	"time"

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
	gwServer := gw_.NewGRPCGateWay(":10000")
	gwServer.RegisterGrpcHandler(func(srv *grpc.Server) {
		date.RegisterDateServiceServer(srv, &s)
	})

	/*
		_ = gwServer.RegisterHTTPFunc(
			context.Background(),
			func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc_.DialOption) error {
				return date.RegisterDateServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
			},
		)
	*/
}
