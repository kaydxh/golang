package stdout

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
)

type StdoutExporterBuilder struct {
}

func NewStdoutExporterBuilder(opts ...StdoutExporterBuilderOption) *StdoutExporterBuilder {

	builder := &StdoutExporterBuilder{}
	builder.ApplyOptions(opts...)
	return builder
}

func (p *StdoutExporterBuilder) Build(
	ctx context.Context,
	c *controller.Controller,
) (aggregation.TemporalitySelector, error) {

	exporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	if err != nil {
		return nil, fmt.Errorf("creating stdoutmetric exporter: %w", err)
	}

	return exporter, nil
}
