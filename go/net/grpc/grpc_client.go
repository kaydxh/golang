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

// connPoolEntry 连接池条目
type connPoolEntry struct {
	conn *grpc.ClientConn
	mu   sync.Mutex
}

var (
	// connPool 全局连接池，key 为 address，value 为 *connPoolEntry
	connPool sync.Map
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
// 对于相同的地址和配置，会复用已存在的连接；如果连接不存在或已关闭，则创建新连接
func GetGrpcClientConn(addr string, disablePrintMethods ...string) (*grpc.ClientConn, error) {
	// 获取或创建连接池条目
	value, _ := connPool.LoadOrStore(addr, &connPoolEntry{})
	entry := value.(*connPoolEntry)

	entry.mu.Lock()
	defer entry.mu.Unlock()

	// 检查现有连接是否可用
	if entry.conn != nil {
		state := entry.conn.GetState()
		// 检查连接是否处于可用状态或可恢复状态
		if isConnAvailable(state) {
			// 检查配置是否一致
			return entry.conn, nil

		} else {
			entry.conn.Close()
			// 连接不可用，清理
			entry.conn = nil
		}
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

	// 保存新连接
	entry.conn = conn

	return conn, nil
}

// isConnAvailable 检查连接状态是否可用
func isConnAvailable(state connectivity.State) bool {
	// Ready: 连接就绪
	// Idle: 连接空闲但可用
	// Connecting: 正在连接，可以等待
	// TransientFailure: 暂时失败，不可用
	// Shutdown: 已关闭，不可用
	return state == connectivity.Ready || state == connectivity.Idle || state == connectivity.Connecting
}

// CloseGrpcClientConn 关闭指定地址的 gRPC 连接并从连接池中移除
func CloseGrpcClientConn(addr string) error {
	if value, ok := connPool.LoadAndDelete(addr); ok {
		entry := value.(*connPoolEntry)
		entry.mu.Lock()
		defer entry.mu.Unlock()
		if entry.conn != nil {
			return entry.conn.Close()
		}
	}
	return nil
}

// CloseAllGrpcClientConns 关闭所有连接池中的连接
func CloseAllGrpcClientConns() error {
	var lastErr error
	var keys []interface{}

	// 先收集所有的 key
	connPool.Range(func(key, value interface{}) bool {
		keys = append(keys, key)
		return true
	})

	// 逐个关闭并删除
	for _, key := range keys {
		if value, ok := connPool.LoadAndDelete(key); ok {
			entry := value.(*connPoolEntry)
			entry.mu.Lock()
			if entry.conn != nil {
				if err := entry.conn.Close(); err != nil {
					lastErr = err
				}
			}
			entry.mu.Unlock()
		}
	}

	return lastErr
}
