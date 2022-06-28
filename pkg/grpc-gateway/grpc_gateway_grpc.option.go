package grpcgateway

import (
	"runtime/debug"

	interceptorlogrus_ "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	interceptorrecovery_ "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	rate_ "github.com/kaydxh/golang/go/time/rate"
	interceptortcloud_ "github.com/kaydxh/golang/pkg/middleware/api/tcloud/v3.0"
	interceptormonitor_ "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/monitor"
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
		interceptorrecovery_.WithRecoveryHandler(func(p interface{}) (err error) {
			logrus.WithError(
				status.Errorf(codes.Internal, "panic triggered: %v at %v", p, string(debug.Stack())),
			).Errorf(
				"recovered in grpc",
			)
			return status.Errorf(codes.Internal, "panic triggered: %v", p)
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

func WithServerUnaryInterceptorsRequestIdOptions() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(c *GRPCGateway) {
		WithServerUnaryInterceptorsOptions(interceptortcloud_.UnaryServerInterceptorOfRequestId()).apply(c)
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
		WithServerUnaryInterceptorsOptions(interceptormonitor_.UnaryServerInterceptorOfInOutPacket()).apply(c)
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
