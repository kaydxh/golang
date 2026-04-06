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

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// UnaryServerTraceInterceptor returns a unary server interceptor that adds OpenTelemetry tracing
// Note: Uses otel.GetTracerProvider() at request time to get the current global TracerProvider
func UnaryServerTraceInterceptor(opts ...otelgrpc.Option) grpc.UnaryServerInterceptor {
	// Get current global TracerProvider and pass it explicitly
	// This ensures the interceptor uses whatever TracerProvider is set at creation time
	tp := otel.GetTracerProvider()
	logrus.Infof("UnaryServerTraceInterceptor: creating trace interceptor, TracerProvider type=%T", tp)
	defaultOpts := []otelgrpc.Option{
		otelgrpc.WithTracerProvider(tp),
	}
	defaultOpts = append(defaultOpts, opts...)
	return otelgrpc.UnaryServerInterceptor(defaultOpts...)
}

// StreamServerTraceInterceptor returns a stream server interceptor that adds OpenTelemetry tracing
func StreamServerTraceInterceptor(opts ...otelgrpc.Option) grpc.StreamServerInterceptor {
	tp := otel.GetTracerProvider()
	logrus.Infof("StreamServerTraceInterceptor: creating trace interceptor, TracerProvider type=%T", tp)
	defaultOpts := []otelgrpc.Option{
		otelgrpc.WithTracerProvider(tp),
	}
	defaultOpts = append(defaultOpts, opts...)
	return otelgrpc.StreamServerInterceptor(defaultOpts...)
}

// UnaryServerTraceInterceptorWithTracer returns a unary server interceptor with custom tracer provider
func UnaryServerTraceInterceptorWithTracer(tracerProvider trace.TracerProvider, opts ...otelgrpc.Option) grpc.UnaryServerInterceptor {
	defaultOpts := []otelgrpc.Option{
		otelgrpc.WithTracerProvider(tracerProvider),
	}
	defaultOpts = append(defaultOpts, opts...)
	return otelgrpc.UnaryServerInterceptor(defaultOpts...)
}

// StreamServerTraceInterceptorWithTracer returns a stream server interceptor with custom tracer provider
func StreamServerTraceInterceptorWithTracer(tracerProvider trace.TracerProvider, opts ...otelgrpc.Option) grpc.StreamServerInterceptor {
	defaultOpts := []otelgrpc.Option{
		otelgrpc.WithTracerProvider(tracerProvider),
	}
	defaultOpts = append(defaultOpts, opts...)
	return otelgrpc.StreamServerInterceptor(defaultOpts...)
}

// GetSpanFromContext extracts the current span from the context
func GetSpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// GetTracerProvider returns the global tracer provider
func GetTracerProvider() trace.TracerProvider {
	return otel.GetTracerProvider()
}

// GetTracer returns a tracer with the given name
func GetTracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return otel.Tracer(name, opts...)
}
