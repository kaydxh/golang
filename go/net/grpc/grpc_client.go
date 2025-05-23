/*
*Copyright (c) 2022, kaydxh
*
*Permission is hereby granted, free of charge, to any person obtaining a copy
go/net/grpc/grpc_client.go*of this software and associated documentation files (the "Software"), to deal
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
package grpc

import (
	"context"
	"fmt"
	"math"
	"time"

	context_ "github.com/kaydxh/golang/go/context"
	interceptordebug_ "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/debug"
	interceptortimer_ "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/timer"
	"google.golang.org/grpc"
)

// 16MB
const (
	defaultMaxMsgSize                      = math.MaxInt32 //16 * 1024 * 1024
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

	ctx, cancel := context_.WithTimeout(context.Background(), c.opts.connectionTimeout)
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

func ClientDialOptions(disablePrintMethods ...string) []grpc.DialOption {
	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(defaultMaxMsgSize),
			grpc.MaxCallSendMsgSize(defaultMaxMsgSize),
		),
		grpc.WithInitialWindowSize(defaultMaxMsgSize),
		grpc.WithInitialConnWindowSize(defaultMaxMsgSize),
		grpc.WithStatsHandler(&statHandler{}),
		grpc.WithChainUnaryInterceptor(
			interceptortimer_.UnaryClientInterceptorOfTimer(),
			interceptordebug_.UnaryClientInterceptorOfInOutputPrinter(disablePrintMethods...)),
	)
	return opts
}

func GetGrpcClientConn(
	addr string,
	connectionTimeout time.Duration,
	disablePrintMethods ...string,
) (*grpc.ClientConn, error) {
	ctx, cancel := context_.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, ClientDialOptions(disablePrintMethods...)...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect address: %v", addr)
	}

	return conn, nil
}
