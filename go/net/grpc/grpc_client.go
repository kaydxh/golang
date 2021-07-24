package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

//16MB
const (
	defaultMaxMsgSize                      = 16 * 1024 * 1024
	defaultConnectionTimeout time.Duration = 5 * time.Second
)

type GrpcClient struct {
	conn *grpc.ClientConn
	opts struct {
		maxMsgSize        int
		connectionTimeout time.Duration
	}
}

func NewGrpcClient(addr string, options ...GrpcClientOption) (*GrpcClient, error) {
	c := &GrpcClient{}
	c.ApplyOptions(options...)

	if c.opts.maxMsgSize == 0 {
		c.opts.maxMsgSize = defaultMaxMsgSize
	}
	if c.opts.connectionTimeout == 0 {
		c.opts.connectionTimeout = defaultConnectionTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.opts.connectionTimeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, ClientDialOptions()...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect address: %v", addr)
	}
	c.conn = conn

	return c, nil
}

func (g *GrpcClient) Conn() *grpc.ClientConn {
	return g.conn
}

func (g *GrpcClient) Close() error {
	if g.conn != nil {
		return g.conn.Close()
	}

	return nil
}

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
