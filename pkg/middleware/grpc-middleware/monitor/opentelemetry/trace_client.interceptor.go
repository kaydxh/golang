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
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// UnaryClientTraceInterceptor returns a unary client interceptor that adds OpenTelemetry tracing
func UnaryClientTraceInterceptor(opts ...otelgrpc.Option) grpc.UnaryClientInterceptor {
	return otelgrpc.UnaryClientInterceptor(opts...)
}

// StreamClientTraceInterceptor returns a stream client interceptor that adds OpenTelemetry tracing
func StreamClientTraceInterceptor(opts ...otelgrpc.Option) grpc.StreamClientInterceptor {
	return otelgrpc.StreamClientInterceptor(opts...)
}

// UnaryClientTraceInterceptorWithTracer returns a unary client interceptor with custom tracer provider
func UnaryClientTraceInterceptorWithTracer(tracerProvider trace.TracerProvider, opts ...otelgrpc.Option) grpc.UnaryClientInterceptor {
	defaultOpts := []otelgrpc.Option{
		otelgrpc.WithTracerProvider(tracerProvider),
	}
	defaultOpts = append(defaultOpts, opts...)
	return otelgrpc.UnaryClientInterceptor(defaultOpts...)
}

// StreamClientTraceInterceptorWithTracer returns a stream client interceptor with custom tracer provider
func StreamClientTraceInterceptorWithTracer(tracerProvider trace.TracerProvider, opts ...otelgrpc.Option) grpc.StreamClientInterceptor {
	defaultOpts := []otelgrpc.Option{
		otelgrpc.WithTracerProvider(tracerProvider),
	}
	defaultOpts = append(defaultOpts, opts...)
	return otelgrpc.StreamClientInterceptor(defaultOpts...)
}

