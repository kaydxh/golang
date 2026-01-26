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

type OTLPTraceExporterBuilderOption interface {
	apply(*OTLPTraceExporterBuilder)
}

type funcOTLPTraceExporterBuilderOption struct {
	f func(*OTLPTraceExporterBuilder)
}

func (fdo *funcOTLPTraceExporterBuilderOption) apply(do *OTLPTraceExporterBuilder) {
	fdo.f(do)
}

func newFuncOTLPTraceExporterBuilderOption(f func(*OTLPTraceExporterBuilder)) *funcOTLPTraceExporterBuilderOption {
	return &funcOTLPTraceExporterBuilderOption{f: f}
}

func (o *OTLPTraceExporterBuilder) ApplyOptions(opts ...OTLPTraceExporterBuilderOption) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

// WithEndpoint sets the target endpoint URL
// For HTTP: include scheme, e.g., "https://tracing.example.com"
// For gRPC: without scheme, e.g., "tracing.example.com:4317"
func WithEndpoint(endpoint string) OTLPTraceExporterBuilderOption {
	return newFuncOTLPTraceExporterBuilderOption(func(o *OTLPTraceExporterBuilder) {
		o.opts.Endpoint = endpoint
	})
}

// WithHeaders sets HTTP headers for authentication
func WithHeaders(headers map[string]string) OTLPTraceExporterBuilderOption {
	return newFuncOTLPTraceExporterBuilderOption(func(o *OTLPTraceExporterBuilder) {
		o.opts.Headers = headers
	})
}

// WithProtocol sets the protocol (HTTP or gRPC)
func WithProtocol(protocol Protocol) OTLPTraceExporterBuilderOption {
	return newFuncOTLPTraceExporterBuilderOption(func(o *OTLPTraceExporterBuilder) {
		o.opts.Protocol = protocol
	})
}

// WithInsecure disables TLS
func WithInsecure(insecure bool) OTLPTraceExporterBuilderOption {
	return newFuncOTLPTraceExporterBuilderOption(func(o *OTLPTraceExporterBuilder) {
		o.opts.Insecure = insecure
	})
}

// WithTimeout sets the timeout for the exporter
func WithTimeout(timeout time.Duration) OTLPTraceExporterBuilderOption {
	return newFuncOTLPTraceExporterBuilderOption(func(o *OTLPTraceExporterBuilder) {
		o.opts.Timeout = timeout
	})
}

// WithURLPath sets the URL path for HTTP protocol
func WithURLPath(urlPath string) OTLPTraceExporterBuilderOption {
	return newFuncOTLPTraceExporterBuilderOption(func(o *OTLPTraceExporterBuilder) {
		o.opts.URLPath = urlPath
	})
}

// WithCompression enables gzip compression
func WithCompression(compression bool) OTLPTraceExporterBuilderOption {
	return newFuncOTLPTraceExporterBuilderOption(func(o *OTLPTraceExporterBuilder) {
		o.opts.Compression = compression
	})
}
