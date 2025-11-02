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
package grpc

import (
	"fmt"
	"math"
	"sync"
	"time"

	interceptordebug_ "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/debug"
	interceptortimer_ "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/timer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	// defaultMaxMsgSize 默认最大消息大小
	defaultMaxMsgSize = math.MaxInt32 // 16 * 1024 * 1024
	// defaultCallTimeout 默认调用超时时间
	defaultCallTimeout = 3 * time.Second
	// defaultKeepaliveTime 默认 keepalive 时间间隔
	defaultKeepaliveTime = 10 * time.Second
	// defaultKeepaliveTimeout 默认 keepalive 超时时间
	defaultKeepaliveTimeout = 3 * time.Second
)

var (
	// connPool 全局连接池，key 为 address，value 为 *grpc.ClientConn
	connPool sync.Map
	// connPoolMutex 连接池互斥锁，用于保证同一地址的连接创建是原子的
	connPoolMutex sync.Map
)

// GrpcClient gRPC 客户端封装
type GrpcClient struct {
	conn *grpc.ClientConn
	opts grpcClientOptions
}

// grpcClientOptions gRPC 客户端配置选项
type grpcClientOptions struct {
	maxMsgSize          int           // 最大消息大小
	disablePrintMethods []string      // 禁止打印的方法列表
	callTimeout         time.Duration // 默认调用超时时间
	keepaliveTime       time.Duration // keepalive 时间间隔
	keepaliveTimeout    time.Duration // keepalive 超时时间
}

// NewGrpcClient 创建一个新的 gRPC 客户端
func NewGrpcClient(addr string, options ...GrpcClientOption) (*GrpcClient, error) {
	c := &GrpcClient{}
	c.ApplyOptions(options...)
	c.setDefaults()

	dialOpts := c.buildDialOptions()
	conn, err := grpc.NewClient(addr, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect address %s: %w", addr, err)
	}
	c.conn = conn

	return c, nil
}

// setDefaults 设置默认配置值
func (g *GrpcClient) setDefaults() {
	if g.opts.maxMsgSize == 0 {
		g.opts.maxMsgSize = defaultMaxMsgSize
	}
	if g.opts.callTimeout == 0 {
		g.opts.callTimeout = defaultCallTimeout
	}
	if g.opts.keepaliveTime == 0 {
		g.opts.keepaliveTime = defaultKeepaliveTime
	}
	if g.opts.keepaliveTimeout == 0 {
		g.opts.keepaliveTimeout = defaultKeepaliveTimeout
	}
}

// buildDialOptions 构建拨号选项
func (g *GrpcClient) buildDialOptions() []grpc.DialOption {
	return ClientDialOptions(
		g.opts.maxMsgSize,
		g.opts.keepaliveTime,
		g.opts.keepaliveTimeout,
		g.opts.disablePrintMethods...,
	)
}

// Conn 返回底层的 gRPC 连接
func (g *GrpcClient) Conn() *grpc.ClientConn {
	return g.conn
}

// CallTimeout 返回默认的调用超时时间
func (g *GrpcClient) CallTimeout() time.Duration {
	return g.opts.callTimeout
}

// Close 关闭 gRPC 连接
func (g *GrpcClient) Close() error {
	if g.conn == nil {
		return nil
	}
	return g.conn.Close()
}

// ClientDialOptions 创建 gRPC 拨号选项
func ClientDialOptions(maxMsgSize int, keepaliveTime, keepaliveTimeout time.Duration, disablePrintMethods ...string) []grpc.DialOption {
	// 设置默认值
	if maxMsgSize == 0 {
		maxMsgSize = defaultMaxMsgSize
	}
	if keepaliveTime == 0 {
		keepaliveTime = defaultKeepaliveTime
	}
	if keepaliveTimeout == 0 {
		keepaliveTimeout = defaultKeepaliveTimeout
	}

	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxMsgSize),
			grpc.MaxCallSendMsgSize(maxMsgSize),
		),
		grpc.WithInitialWindowSize(int32(maxMsgSize)),
		grpc.WithInitialConnWindowSize(int32(maxMsgSize)),
		grpc.WithStatsHandler(&statHandler{}),
		grpc.WithChainUnaryInterceptor(
			interceptortimer_.UnaryClientInterceptorOfTimer(),
			interceptordebug_.UnaryClientInterceptorOfInOutputPrinter(disablePrintMethods...),
		),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                keepaliveTime,
			Timeout:             keepaliveTimeout,
			PermitWithoutStream: true,
		}),
	}
}

// GetGrpcClientConn 获取一个 gRPC 客户端长连接（支持连接复用）
// 对于相同的地址，会复用已存在的连接；如果连接不存在或已关闭，则创建新连接
func GetGrpcClientConn(addr string, disablePrintMethods ...string) (*grpc.ClientConn, error) {
	// 尝试从连接池获取已存在的连接
	if conn, ok := connPool.Load(addr); ok {
		clientConn := conn.(*grpc.ClientConn)
		// 检查连接状态，如果连接正常则直接返回
		state := clientConn.GetState()
		if state != connectivity.Shutdown && state != connectivity.TransientFailure {
			return clientConn, nil
		}
		// 连接已关闭或失败，从池中删除
		connPool.Delete(addr)
	}

	// 使用互斥锁确保同一地址只创建一次连接
	mutex, _ := connPoolMutex.LoadOrStore(addr, &sync.Mutex{})
	mu := mutex.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	// 双重检查：可能在等待锁的过程中，其他 goroutine 已经创建了连接
	if conn, ok := connPool.Load(addr); ok {
		clientConn := conn.(*grpc.ClientConn)
		state := clientConn.GetState()
		if state != connectivity.Shutdown && state != connectivity.TransientFailure {
			return clientConn, nil
		}
		connPool.Delete(addr)
	}

	// 创建新连接
	opts := ClientDialOptions(
		defaultMaxMsgSize,
		defaultKeepaliveTime,
		defaultKeepaliveTimeout,
		disablePrintMethods...,
	)
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client for address %s: %w", addr, err)
	}

	// 将新连接存入连接池
	connPool.Store(addr, conn)

	return conn, nil
}

// CloseGrpcClientConn 关闭指定地址的 gRPC 连接并从连接池中移除
func CloseGrpcClientConn(addr string) error {
	if conn, ok := connPool.LoadAndDelete(addr); ok {
		return conn.(*grpc.ClientConn).Close()
	}
	return nil
}

// CloseAllGrpcClientConns 关闭所有连接池中的连接
func CloseAllGrpcClientConns() error {
	var lastErr error
	connPool.Range(func(key, value interface{}) bool {
		if err := value.(*grpc.ClientConn).Close(); err != nil {
			lastErr = err
		}
		connPool.Delete(key)
		return true
	})
	return lastErr
}
