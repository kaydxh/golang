package stdout

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

type StdoutExporterBuilderOptions struct {
	prettyPrint bool
}

type StdoutExporterBuilder struct {
	opts StdoutExporterBuilderOptions
}

func defaultBuilderOptions() StdoutExporterBuilderOptions {
	return StdoutExporterBuilderOptions{
		prettyPrint: true,
	}
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

	var opts []stdoutmetric.Option
	if p.opts.prettyPrint {
		opts = append(opts, stdoutmetric.WithPrettyPrint())
	}

	exporter, err := stdoutmetric.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("creating stdoutmetric exporter: %w", err)
	}

	return exporter, nil
}
