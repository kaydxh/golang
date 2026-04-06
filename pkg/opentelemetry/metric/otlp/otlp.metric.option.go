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

import "time"

type OTLPExporterBuilderOption interface {
	apply(*OTLPExporterBuilder)
}

type funcOTLPExporterBuilderOption struct {
	f func(*OTLPExporterBuilder)
}

func (fdo *funcOTLPExporterBuilderOption) apply(do *OTLPExporterBuilder) {
	fdo.f(do)
}

func newFuncOTLPExporterBuilderOption(f func(*OTLPExporterBuilder)) *funcOTLPExporterBuilderOption {
	return &funcOTLPExporterBuilderOption{f: f}
}

func (o *OTLPExporterBuilder) ApplyOptions(opts ...OTLPExporterBuilderOption) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

// WithEndpoint sets the target endpoint URL
// For HTTP: include scheme, e.g., "https://prometheus.example.com"
// For gRPC: without scheme, e.g., "prometheus.example.com:4317"
func WithEndpoint(endpoint string) OTLPExporterBuilderOption {
	return newFuncOTLPExporterBuilderOption(func(o *OTLPExporterBuilder) {
		o.opts.Endpoint = endpoint
	})
}

// WithHeaders sets HTTP headers for authentication
// For Tencent Cloud Prometheus, set Authorization header with Token
func WithHeaders(headers map[string]string) OTLPExporterBuilderOption {
	return newFuncOTLPExporterBuilderOption(func(o *OTLPExporterBuilder) {
		o.opts.Headers = headers
	})
}

// WithProtocol sets the protocol (HTTP or gRPC)
func WithProtocol(protocol Protocol) OTLPExporterBuilderOption {
	return newFuncOTLPExporterBuilderOption(func(o *OTLPExporterBuilder) {
		o.opts.Protocol = protocol
	})
}

// WithInsecure disables TLS
func WithInsecure(insecure bool) OTLPExporterBuilderOption {
	return newFuncOTLPExporterBuilderOption(func(o *OTLPExporterBuilder) {
		o.opts.Insecure = insecure
	})
}

// WithTimeout sets the timeout for the exporter
func WithTimeout(timeout time.Duration) OTLPExporterBuilderOption {
	return newFuncOTLPExporterBuilderOption(func(o *OTLPExporterBuilder) {
		o.opts.Timeout = timeout
	})
}

// WithURLPath sets the URL path for HTTP protocol
func WithURLPath(urlPath string) OTLPExporterBuilderOption {
	return newFuncOTLPExporterBuilderOption(func(o *OTLPExporterBuilder) {
		o.opts.URLPath = urlPath
	})
}

// WithCompression enables gzip compression
func WithCompression(compression bool) OTLPExporterBuilderOption {
	return newFuncOTLPExporterBuilderOption(func(o *OTLPExporterBuilder) {
		o.opts.Compression = compression
	})
}

// WithTemporalityDelta uses Delta temporality instead of Cumulative
// ZhiYan platform requires Delta temporality
func WithTemporalityDelta(delta bool) OTLPExporterBuilderOption {
	return newFuncOTLPExporterBuilderOption(func(o *OTLPExporterBuilder) {
		o.opts.TemporalityDelta = delta
	})
}
