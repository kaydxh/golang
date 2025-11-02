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

import "time"

// WithMaxMsgSize 设置最大消息大小
func WithMaxMsgSize(maxMsgSize int) GrpcClientOption {
	return GrpcClientOptionFunc(func(c *GrpcClient) {
		c.opts.maxMsgSize = maxMsgSize
	})
}

// WithCallTimeout 设置默认调用超时时间
func WithCallTimeout(timeout time.Duration) GrpcClientOption {
	return GrpcClientOptionFunc(func(c *GrpcClient) {
		c.opts.callTimeout = timeout
	})
}

// WithKeepaliveTime 设置 keepalive 时间间隔
func WithKeepaliveTime(keepaliveTime time.Duration) GrpcClientOption {
	return GrpcClientOptionFunc(func(c *GrpcClient) {
		c.opts.keepaliveTime = keepaliveTime
	})
}

// WithKeepaliveTimeout 设置 keepalive 超时时间
func WithKeepaliveTimeout(keepaliveTimeout time.Duration) GrpcClientOption {
	return GrpcClientOptionFunc(func(c *GrpcClient) {
		c.opts.keepaliveTimeout = keepaliveTimeout
	})
}

// WithDisablePrintMethods 设置不需要在调试日志中打印的方法列表
func WithDisablePrintMethods(methods ...string) GrpcClientOption {
	return GrpcClientOptionFunc(func(c *GrpcClient) {
		c.opts.disablePrintMethods = methods
	})
}
