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
