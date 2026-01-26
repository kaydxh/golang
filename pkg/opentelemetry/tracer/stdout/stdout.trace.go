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
package stdout

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type StdoutExporterBuilderOptions struct {
	stdouttraceOpts []stdouttrace.Option
}

type StdoutExporterBuilder struct {
	opts StdoutExporterBuilderOptions
}

func defaultBuilderOptions() StdoutExporterBuilderOptions {
	return StdoutExporterBuilderOptions{}
}

func NewStdoutExporterBuilder(opts ...StdoutExporterBuilderOption) *StdoutExporterBuilder {

	builder := &StdoutExporterBuilder{
		opts: defaultBuilderOptions(),
	}
	builder.ApplyOptions(opts...)
	return builder
}

func (p *StdoutExporterBuilder) Build(
	ctx context.Context,
) (sdktrace.SpanExporter, error) {

	exporter, err := stdouttrace.New(p.opts.stdouttraceOpts...)
	if err != nil {
		return nil, fmt.Errorf("creating stdouttrace exporter: %w", err)
	}

	return exporter, nil
}
