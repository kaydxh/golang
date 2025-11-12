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
package interceptoropentelemetry

// Example usage:
//
// Server side:
//
//	import (
//		"context"
//		interceptoropentelemetry "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/monitor/opentelemetry"
//		"github.com/kaydxh/golang/pkg/monitor/opentelemetry/tracer"
//		"github.com/kaydxh/golang/pkg/monitor/opentelemetry/tracer/jaeger"
//		"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
//		"google.golang.org/grpc"
//	)
//
//	func createGrpcServer() (*grpc.Server, error) {
//		// Initialize tracer (optional, uses global tracer if not set)
//		jaegerBuilder, err := jaeger.NewJaegerExporertBuilder("http://localhost:14268/api/traces")
//		if err != nil {
//			return nil, err
//		}
//
//		t := tracer.NewTracer(
//			tracer.WithExporterBuilder(jaegerBuilder),
//			tracer.WithServiceName("my-grpc-service"),
//			tracer.WithServiceVersion("1.0.0"),
//		)
//		if err := t.Install(context.Background()); err != nil {
//			return nil, err
//		}
//
//		// Create gRPC server with OpenTelemetry interceptors
//		server := grpc.NewServer(
//			grpc.ChainUnaryInterceptor(
//				interceptoropentelemetry.UnaryServerTraceInterceptor(),
//				// Add other interceptors as needed
//			),
//			grpc.ChainStreamInterceptor(
//				interceptoropentelemetry.StreamServerTraceInterceptor(),
//			),
//		)
//
//		return server, nil
//	}
//
// Client side:
//
//	import (
//		interceptoropentelemetry "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/monitor/opentelemetry"
//		"google.golang.org/grpc"
//	)
//
//	func createGrpcClient(addr string) (*grpc.ClientConn, error) {
//		// Create gRPC client with OpenTelemetry interceptors
//		conn, err := grpc.NewClient(addr,
//			grpc.WithChainUnaryInterceptor(
//				interceptoropentelemetry.UnaryClientTraceInterceptor(),
//				// Add other interceptors as needed
//			),
//			grpc.WithChainStreamInterceptor(
//				interceptoropentelemetry.StreamClientTraceInterceptor(),
//			),
//		)
//		if err != nil {
//			return nil, err
//		}
//
//		return conn, nil
//	}
//
// Using custom tracer provider:
//
//	import (
//		interceptoropentelemetry "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/monitor/opentelemetry"
//		"github.com/kaydxh/golang/pkg/monitor/opentelemetry/tracer"
//		"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
//		"google.golang.org/grpc"
//	)
//
//	func createGrpcServerWithCustomTracer() (*grpc.Server, error) {
//		// Create and install tracer
//		t := tracer.NewTracer(
//			tracer.WithExporterBuilder(myExporterBuilder),
//			tracer.WithServiceName("my-service"),
//		)
//		if err := t.Install(context.Background()); err != nil {
//			return nil, err
//		}
//
//		// Create server with custom tracer provider
//		server := grpc.NewServer(
//			grpc.ChainUnaryInterceptor(
//				interceptoropentelemetry.UnaryServerTraceInterceptorWithTracer(
//					t.TracerProvider(),
//					otelgrpc.WithMessageEvents(otelgrpc.SentEvents, otelgrpc.ReceivedEvents),
//				),
//			),
//		)
//
//		return server, nil
//	}
//
// Extracting span from context:
//
//	import (
//		interceptoropentelemetry "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/monitor/opentelemetry"
//		"go.opentelemetry.io/otel/attribute"
//		"go.opentelemetry.io/otel/codes"
//	)
//
//	func myGrpcHandler(ctx context.Context, req *MyRequest) (*MyResponse, error) {
//		// Get the current span from context
//		span := interceptoropentelemetry.GetSpanFromContext(ctx)
//
//		// Add custom attributes
//		span.SetAttributes(
//			attribute.String("custom.key", "value"),
//			attribute.Int("user.id", 123),
//		)
//
//		// Add events
//		span.AddEvent("processing request")
//
//		// Set status
//		span.SetStatus(codes.Ok, "success")
//
//		// ... your handler logic ...
//
//		return &MyResponse{}, nil
//	}

