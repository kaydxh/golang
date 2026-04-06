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

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Protocol int

const (
	ProtocolHTTP Protocol = iota
	ProtocolGRPC
)

type OTLPExporterBuilderOptions struct {
	// Endpoint is the target endpoint URL (without scheme for gRPC, with scheme for HTTP)
	// For Tencent Cloud Prometheus, use the Remote Write address
	Endpoint string

	// Headers for authentication (e.g., Token for Tencent Cloud Prometheus)
	Headers map[string]string

	// Protocol specifies HTTP or gRPC
	Protocol Protocol

	// Insecure disables TLS (default: false)
	Insecure bool

	// Timeout for the exporter
	Timeout time.Duration

	// URLPath is the URL path for HTTP protocol (default: "/v1/metrics")
	URLPath string

	// Compression enables gzip compression
	Compression bool

	// TemporalityDelta uses Delta temporality instead of Cumulative
	TemporalityDelta bool
}

type OTLPExporterBuilder struct {
	opts OTLPExporterBuilderOptions
}

func defaultBuilderOptions() OTLPExporterBuilderOptions {
	return OTLPExporterBuilderOptions{
		Protocol: ProtocolHTTP,
		Insecure: false,
		Timeout:  30 * time.Second,
		URLPath:  "/v1/metrics",
	}
}

func NewOTLPExporterBuilder(opts ...OTLPExporterBuilderOption) *OTLPExporterBuilder {
	builder := &OTLPExporterBuilder{
		opts: defaultBuilderOptions(),
	}
	builder.ApplyOptions(opts...)
	return builder
}

func (p *OTLPExporterBuilder) Build(ctx context.Context) (metric.Exporter, error) {
	switch p.opts.Protocol {
	case ProtocolHTTP:
		return p.buildHTTPExporter(ctx)
	case ProtocolGRPC:
		return p.buildGRPCExporter(ctx)
	default:
		return nil, fmt.Errorf("unsupported protocol: %d", p.opts.Protocol)
	}
}

func (p *OTLPExporterBuilder) buildHTTPExporter(ctx context.Context) (metric.Exporter, error) {
	if p.opts.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required for OTLP HTTP exporter")
	}

	opts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(p.opts.Endpoint),
	}

	if p.opts.Insecure {
		opts = append(opts, otlpmetrichttp.WithInsecure())
	}

	if len(p.opts.Headers) > 0 {
		opts = append(opts, otlpmetrichttp.WithHeaders(p.opts.Headers))
	}

	if p.opts.Timeout > 0 {
		opts = append(opts, otlpmetrichttp.WithTimeout(p.opts.Timeout))
	}

	if p.opts.URLPath != "" {
		opts = append(opts, otlpmetrichttp.WithURLPath(p.opts.URLPath))
	}

	if p.opts.Compression {
		opts = append(opts, otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression))
	}

	if p.opts.TemporalityDelta {
		opts = append(opts, otlpmetrichttp.WithTemporalitySelector(deltaTemporalitySelector))
	}

	exporter, err := otlpmetrichttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP HTTP metric exporter: %w", err)
	}

	return exporter, nil
}

func (p *OTLPExporterBuilder) buildGRPCExporter(ctx context.Context) (metric.Exporter, error) {
	if p.opts.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required for OTLP gRPC exporter")
	}

	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(p.opts.Endpoint),
	}

	if p.opts.Insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
		opts = append(opts, otlpmetricgrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
	}

	if len(p.opts.Headers) > 0 {
		opts = append(opts, otlpmetricgrpc.WithHeaders(p.opts.Headers))
	}

	if p.opts.Timeout > 0 {
		opts = append(opts, otlpmetricgrpc.WithTimeout(p.opts.Timeout))
	}

	if p.opts.Compression {
		opts = append(opts, otlpmetricgrpc.WithCompressor("gzip"))
	}

	if p.opts.TemporalityDelta {
		opts = append(opts, otlpmetricgrpc.WithTemporalitySelector(deltaTemporalitySelector))
	}

	exporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP gRPC metric exporter: %w", err)
	}

	return exporter, nil
}

// deltaTemporalitySelector returns Delta temporality for all instrument kinds
// ZhiYan platform requires Delta temporality
func deltaTemporalitySelector(kind metric.InstrumentKind) metricdata.Temporality {
	return metricdata.DeltaTemporality
}
