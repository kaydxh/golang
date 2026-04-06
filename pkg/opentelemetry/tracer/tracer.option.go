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
package tracer

import (
	"go.opentelemetry.io/otel/sdk/resource"
)

func WithExporterBuilder(builder TracerExporterBuilder) TracerOption {
	return TracerOptionFunc(func(m *Tracer) {
		m.opts.builer = builder
	})
}

// WithServiceName sets the service name for the tracer
func WithServiceName(name string) TracerOption {
	return TracerOptionFunc(func(m *Tracer) {
		m.opts.serviceName = name
	})
}

// WithServiceVersion sets the service version for the tracer
func WithServiceVersion(version string) TracerOption {
	return TracerOptionFunc(func(m *Tracer) {
		m.opts.serviceVersion = version
	})
}

// WithServiceNamespace sets the service namespace for the tracer
func WithServiceNamespace(namespace string) TracerOption {
	return TracerOptionFunc(func(m *Tracer) {
		m.opts.serviceNamespace = namespace
	})
}

// WithResource sets a custom resource for the tracer
func WithResource(res *resource.Resource) TracerOption {
	return TracerOptionFunc(func(m *Tracer) {
		m.opts.resource = res
	})
}

// WithExporterLogging enables logging for trace export operations
// name: exporter name for logging (e.g., "OTLP", "Jaeger", "stdout")
// endpoint: endpoint address for logging
func WithExporterLogging(name, endpoint string) TracerOption {
	return TracerOptionFunc(func(m *Tracer) {
		m.opts.exporterName = name
		m.opts.exporterEndpoint = endpoint
	})
}
