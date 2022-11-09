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
package grpcgateway

import (
	"runtime/debug"
	"time"

	interceptorlogrus_ "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	interceptorrecovery_ "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	runtime_ "github.com/kaydxh/golang/go/runtime"
	rate_ "github.com/kaydxh/golang/go/time/rate"
	interceptortcloud_ "github.com/kaydxh/golang/pkg/middleware/api/tcloud/v3.0"
	interceptordebug_ "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/debug"
	interceptoropentelemetry_ "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/monitor/opentelemetry"
	interceptorprometheus_ "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/monitor/prometheus"
	interceptorratelimit_ "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/ratelimit"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WithServerOptions(opts ...grpc.ServerOption) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.serverOptions = append(c.opts.serverOptions, opts...)
	})
}

func WithClientDialOptions(opts ...grpc.DialOption) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.clientDialOptions = append(c.opts.clientDialOptions, opts...)
	})
}

func WithServerUnaryInterceptorsOptions(opts ...grpc.UnaryServerInterceptor) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.interceptionOptions.grpcServerOpts.unaryInterceptors = append(
			c.opts.interceptionOptions.grpcServerOpts.unaryInterceptors,
			opts...)
	})
}

func WithServerStreamInterceptorsOptions(opts ...grpc.StreamServerInterceptor) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.interceptionOptions.grpcServerOpts.streamInterceptors = append(
			c.opts.interceptionOptions.grpcServerOpts.streamInterceptors,
			opts...)
	})
}

func WithServerInterceptorsLogrusOptions(
	logger *logrus.Logger,
) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		l := logrus.NewEntry(logger)
		WithServerUnaryInterceptorsOptions(interceptorlogrus_.UnaryServerInterceptor(l)).apply(c)
		WithServerStreamInterceptorsOptions(interceptorlogrus_.StreamServerInterceptor(l)).apply(c)
	})
}

// recover
func WithServerInterceptorsRecoveryOptions() GRPCGatewayOption {
	opts := []interceptorrecovery_.Option{
		interceptorrecovery_.WithRecoveryHandler(func(r interface{}) (err error) {
			out, err := runtime_.FormatStack()
			if err != nil {
				logrus.WithError(status.Errorf(codes.Internal, "%s", r)).Errorf("%s", debug.Stack())
			} else {
				logrus.WithError(status.Errorf(codes.Internal, "%s", r)).Errorf("%s", out)
			}

			return status.Errorf(codes.Internal, "panic: %v", r)
		}),
	}

	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(interceptorrecovery_.UnaryServerInterceptor(opts...)).apply(c)
		WithServerStreamInterceptorsOptions(interceptorrecovery_.StreamServerInterceptor(opts...)).apply(c)
	})
}

func WithServerUnaryInterceptorsTimerOptions(enabledMetric bool) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(interceptorprometheus_.UnaryServerInterceptorOfTimer(enabledMetric)).apply(c)
	})
}

func WithServerUnaryInterceptorsCodeMessageOptions(enabledMetric bool) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(
			interceptorprometheus_.UnaryServerInterceptorOfCodeMessage(enabledMetric),
		).apply(
			c,
		)
	})
}

func WithServerUnaryMetricInterceptorOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(
			interceptoropentelemetry_.UnaryServerMetricInterceptor(),
		).apply(
			c,
		)
	})
}

func WithServerUnaryInterceptorsRequestIdOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(interceptordebug_.UnaryServerInterceptorOfRequestId()).apply(c)
	})
}

func WithServerUnaryInterceptorsErrorOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(interceptortcloud_.UnaryServerInterceptorOfError()).apply(c)
		//WithServerStreamInterceptorsOptions(interceptortcloud_.StreamServerInterceptor()).apply(c)
	})
}

//limiter rate for grpc api
func WithServerInterceptorsLimitRateOptions(burstUnary, burstStream int) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		if burstUnary > 0 {
			limiterUnary := rate_.NewLimiter(burstUnary)
			WithServerUnaryInterceptorsOptions(interceptorratelimit_.UnaryServerInterceptor(limiterUnary)).apply(c)
		}

		if burstStream > 0 {
			limiterStream := rate_.NewLimiter(burstStream)
			WithServerStreamInterceptorsOptions(interceptorratelimit_.StreamServerInterceptor(limiterStream)).apply(c)
		}
	})
}

func WithServerUnaryInterceptorsInOutPacketOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(interceptordebug_.UnaryServerInterceptorOfInOutputPrinter()).apply(c)
	})
}

// timeout
func WithServerInterceptorTimeoutOptions(timeout time.Duration) GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		if timeout > 0 {
			runtime.DefaultContextTimeout = timeout
		}
	})
}

/*
// WithHttpHandlerInterceptorOptions
func WithHttpHandlerInterceptorOptions(opts ...http_.HandlerChainOption) GRPCGatewayOption {

	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.interceptionOptions.httpServerOpts.handlerChain.ApplyOptions(opts...)
	})
}

// WithHttpHandlerInterceptorOptions
func WithHttpHandlerInterceptorOptions(opts ...http_.HandlerInterceptorsOption) GRPCGatewayOption {

	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		c.opts.interceptionOptions.httpServerOpts.httpInterceptors.ApplyOptions(opts...)
	})
}

func WithHttpHandlerInterceptor(handler func(http.Handler) http.Handler) GRPCGatewayOption {

	return WithHttpHandlerInterceptorOptions(
		http_.HandlerInterceptorsOptionFunc(func(handlers *http_.HandlerInterceptors) {
			handlers.Interceptors = append(handlers.Interceptors, handler)
		}),
	)
}
*/

// WithHttpHandlerInterceptorTraceIDOptions
/*
func WithHttpHandlerInterceptorInOutPacketOptions() GRPCGatewayOption {
	//return WithHttpHandlerInterceptor(httpinterceptorinoutpacket_.InOutPacket)
	return WithHttpHandlerInterceptor(httpinterceptorinoutpacket_.InOutPacket)
}
*/
