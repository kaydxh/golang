package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

//16MB
const defaultMaxMsgSize = 16 * 1024 * 1024

func ClientDialOptions() []grpc.DialOption {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(defaultMaxMsgSize)))
	return opts
}

func GetGrpcClientConn(addr string, connectionTimeout time.Duration) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, ClientDialOptions()...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect address: %v", addr)
	}

	return conn, nil
}
