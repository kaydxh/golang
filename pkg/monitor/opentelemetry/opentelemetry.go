package opentelemetry

import (
	"context"

	metric_ "github.com/kaydxh/golang/pkg/monitor/opentelemetry/metric"
	tracer_ "github.com/kaydxh/golang/pkg/monitor/opentelemetry/tracer"
)

type OpenTelemetryOptions struct {
	meterOptions  []metric_.MeterOption
	tracerOptions []tracer_.TracerOption
}

type OpenTelemetry struct {
	opts OpenTelemetryOptions
}

func NewOpenTelemetry(opts ...OpenTelemetryOption) *OpenTelemetry {
	t := &OpenTelemetry{}
	t.ApplyOptions(opts...)

	return t
}

func (t *OpenTelemetry) Install(ctx context.Context) error {

	if len(t.opts.meterOptions) > 0 {
		meter := metric_.NewMeter(t.opts.meterOptions...)
		err := meter.Install(ctx)
		if err != nil {
			return err
		}
	}

	if len(t.opts.tracerOptions) > 0 {
		tracer := tracer_.NewTracer()
		err := tracer.Install(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
