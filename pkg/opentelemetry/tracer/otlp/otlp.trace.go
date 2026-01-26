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
package otlp

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Protocol int

const (
	ProtocolHTTP Protocol = iota
	ProtocolGRPC
)

type OTLPTraceExporterBuilderOptions struct {
	// Endpoint is the target endpoint URL (without scheme for gRPC, with scheme for HTTP)
	Endpoint string

	// Headers for authentication
	Headers map[string]string

	// Protocol specifies HTTP or gRPC
	Protocol Protocol

	// Insecure disables TLS (default: false)
	Insecure bool

	// Timeout for the exporter
	Timeout time.Duration

	// URLPath is the URL path for HTTP protocol (default: "/v1/traces")
	URLPath string

	// Compression enables gzip compression
	Compression bool
}

type OTLPTraceExporterBuilder struct {
	opts OTLPTraceExporterBuilderOptions
}

func defaultBuilderOptions() OTLPTraceExporterBuilderOptions {
	return OTLPTraceExporterBuilderOptions{
		Protocol: ProtocolHTTP,
		Insecure: false,
		Timeout:  30 * time.Second,
		URLPath:  "/v1/traces",
	}
}

func NewOTLPTraceExporterBuilder(opts ...OTLPTraceExporterBuilderOption) *OTLPTraceExporterBuilder {
	builder := &OTLPTraceExporterBuilder{
		opts: defaultBuilderOptions(),
	}
	builder.ApplyOptions(opts...)
	return builder
}

func (b *OTLPTraceExporterBuilder) Build(ctx context.Context) (sdktrace.SpanExporter, error) {
	switch b.opts.Protocol {
	case ProtocolHTTP:
		return b.buildHTTPExporter(ctx)
	case ProtocolGRPC:
		return b.buildGRPCExporter(ctx)
	default:
		return nil, fmt.Errorf("unsupported protocol: %d", b.opts.Protocol)
	}
}

func (b *OTLPTraceExporterBuilder) buildHTTPExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	if b.opts.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required for OTLP HTTP trace exporter")
	}

	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(b.opts.Endpoint),
	}

	if b.opts.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	if len(b.opts.Headers) > 0 {
		opts = append(opts, otlptracehttp.WithHeaders(b.opts.Headers))
	}

	if b.opts.Timeout > 0 {
		opts = append(opts, otlptracehttp.WithTimeout(b.opts.Timeout))
	}

	if b.opts.URLPath != "" {
		opts = append(opts, otlptracehttp.WithURLPath(b.opts.URLPath))
	}

	if b.opts.Compression {
		opts = append(opts, otlptracehttp.WithCompression(otlptracehttp.GzipCompression))
	}

	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP HTTP trace exporter: %w", err)
	}

	return exporter, nil
}

func (b *OTLPTraceExporterBuilder) buildGRPCExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	if b.opts.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required for OTLP gRPC trace exporter")
	}

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(b.opts.Endpoint),
	}

	if b.opts.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
		opts = append(opts, otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
	}

	if len(b.opts.Headers) > 0 {
		opts = append(opts, otlptracegrpc.WithHeaders(b.opts.Headers))
	}

	if b.opts.Timeout > 0 {
		opts = append(opts, otlptracegrpc.WithTimeout(b.opts.Timeout))
	}

	if b.opts.Compression {
		opts = append(opts, otlptracegrpc.WithCompressor("gzip"))
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP gRPC trace exporter: %w", err)
	}

	return exporter, nil
}
