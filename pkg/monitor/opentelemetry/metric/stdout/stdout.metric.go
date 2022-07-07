package stdout

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

type StdoutExporterBuilderOptions struct {
	stdoutmetricOpts []stdoutmetric.Option
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
) (export.Exporter, error) {

	exporter, err := stdoutmetric.New(p.opts.stdoutmetricOpts...)
	if err != nil {
		return nil, fmt.Errorf("creating stdoutmetric exporter: %w", err)
	}

	return exporter, nil
}
