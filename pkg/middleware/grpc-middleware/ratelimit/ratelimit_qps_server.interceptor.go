/*
 *Copyright (c) 2024, kaydxh
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
package ratelimit

import (
	"context"
	"time"

	rate_ "github.com/kaydxh/golang/go/time/rate"
	logs_ "github.com/kaydxh/golang/pkg/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QPSLimiter is the interface for QPS-based rate limiting.
type QPSLimiter interface {
	Allow(method string) bool
	AllowFor(method string, timeout time.Duration) bool
	Wait(ctx context.Context, method string) error
}

// UnaryServerInterceptorQPS returns a new unary server interceptor
// that performs QPS-based rate limiting with per-method configuration.
func UnaryServerInterceptorQPS(limiter QPSLimiter) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if !limiter.Allow(info.FullMethod) {
			logger := logs_.GetLogger(ctx)
			err := status.Errorf(
				codes.ResourceExhausted,
				"%s is rejected by grpc_ratelimit_qps middleware, QPS limit exceeded, please retry later.",
				info.FullMethod,
			)
			logger.WithField("method", info.FullMethod).Errorf("rate limited: %v", err)
			return nil, err
		}
		return handler(ctx, req)
	}
}

// UnaryServerInterceptorQPSWait returns a new unary server interceptor
// that waits for rate limit token instead of rejecting immediately.
func UnaryServerInterceptorQPSWait(limiter QPSLimiter, timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Try with timeout
		if !limiter.AllowFor(info.FullMethod, timeout) {
			logger := logs_.GetLogger(ctx)
			err := status.Errorf(
				codes.ResourceExhausted,
				"%s is rejected by grpc_ratelimit_qps middleware after waiting %v, QPS limit exceeded.",
				info.FullMethod,
				timeout,
			)
			logger.WithField("method", info.FullMethod).Errorf("rate limited: %v", err)
			return nil, err
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptorQPS returns a new stream server interceptor
// that performs QPS-based rate limiting with per-method configuration.
func StreamServerInterceptorQPS(limiter QPSLimiter) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if !limiter.Allow(info.FullMethod) {
			logger := logs_.GetLogger(stream.Context())
			err := status.Errorf(
				codes.ResourceExhausted,
				"%s is rejected by grpc_ratelimit_qps middleware, QPS limit exceeded, please retry later.",
				info.FullMethod,
			)
			logger.WithField("method", info.FullMethod).Errorf("rate limited: %v", err)
			return err
		}
		return handler(srv, stream)
	}
}

// StreamServerInterceptorQPSWait returns a new stream server interceptor
// that waits for rate limit token instead of rejecting immediately.
func StreamServerInterceptorQPSWait(limiter QPSLimiter, timeout time.Duration) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if !limiter.AllowFor(info.FullMethod, timeout) {
			logger := logs_.GetLogger(stream.Context())
			err := status.Errorf(
				codes.ResourceExhausted,
				"%s is rejected by grpc_ratelimit_qps middleware after waiting %v, QPS limit exceeded.",
				info.FullMethod,
				timeout,
			)
			logger.WithField("method", info.FullMethod).Errorf("rate limited: %v", err)
			return err
		}
		return handler(srv, stream)
	}
}

// NewMethodQPSLimiter creates a pre-configured method QPS limiter.
// This is a convenience wrapper for rate_.NewMethodQPSLimiter.
func NewMethodQPSLimiter(defaultQPS float64, defaultBurst int) *rate_.MethodQPSLimiter {
	return rate_.NewMethodQPSLimiter(defaultQPS, defaultBurst)
}

// NewMethodQPSLimiterWithConfigs creates a limiter with method configurations.
func NewMethodQPSLimiterWithConfigs(
	defaultQPS float64,
	defaultBurst int,
	configs []rate_.MethodQPSConfig,
) (*rate_.MethodQPSLimiter, error) {
	return rate_.NewMethodQPSLimiterWithConfigs(defaultQPS, defaultBurst, configs)
}

// ConcurrencyLimiter is the interface for concurrency-based rate limiting.
// Unlike QPS limiting, concurrency control limits the number of requests being processed simultaneously.
// Tokens are returned when request processing completes.
type ConcurrencyLimiter interface {
	Allow(method string) bool
	Put(method string)
}

// UnaryServerInterceptorConcurrency returns a new unary server interceptor
// that performs concurrency-based rate limiting with per-method configuration.
func UnaryServerInterceptorConcurrency(limiter ConcurrencyLimiter) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if !limiter.Allow(info.FullMethod) {
			logger := logs_.GetLogger(ctx)
			err := status.Errorf(
				codes.ResourceExhausted,
				"%s is rejected by grpc_concurrency middleware, max concurrency exceeded, please retry later.",
				info.FullMethod,
			)
			logger.WithField("method", info.FullMethod).Errorf("concurrency limited: %v", err)
			return nil, err
		}
		defer limiter.Put(info.FullMethod)
		return handler(ctx, req)
	}
}

// StreamServerInterceptorConcurrency returns a new stream server interceptor
// that performs concurrency-based rate limiting with per-method configuration.
func StreamServerInterceptorConcurrency(limiter ConcurrencyLimiter) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if !limiter.Allow(info.FullMethod) {
			logger := logs_.GetLogger(stream.Context())
			err := status.Errorf(
				codes.ResourceExhausted,
				"%s is rejected by grpc_concurrency middleware, max concurrency exceeded, please retry later.",
				info.FullMethod,
			)
			logger.WithField("method", info.FullMethod).Errorf("concurrency limited: %v", err)
			return err
		}
		defer limiter.Put(info.FullMethod)
		return handler(srv, stream)
	}
}
