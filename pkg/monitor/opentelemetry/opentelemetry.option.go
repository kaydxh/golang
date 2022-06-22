package opentelemetry

import (
	"time"

	"github.com/kaydxh/golang/pkg/monitor/opentelemetry/metrics/meter"
	"github.com/kaydxh/golang/pkg/monitor/opentelemetry/tracer"
)

func WithMeterExporter(exporterBuilder meter.ExporterBuilder) OpenTelemetryOption {
	return OpenTelemetryOptionFunc(func(o *OpenTelemetry) {
		o.opts.meterOptions = append(o.opts.meterOptions, meter.WithExporter(exporterBuilder))

	})
}

func WithMeterPullExporter(pullExporterBuilder meter.PullExporterBuilder) OpenTelemetryOption {
	return OpenTelemetryOptionFunc(func(o *OpenTelemetry) {
		o.opts.meterOptions = append(o.opts.meterOptions, meter.WithPullExporter(pullExporterBuilder))
	})
}

func WithMetricCollectDuration(period time.Duration) OpenTelemetryOption {
	return OpenTelemetryOptionFunc(func(o *OpenTelemetry) {
		o.opts.meterOptions = append(o.opts.meterOptions, meter.WithCollectPeriod(period))
	})
}

func WithTracerExporter(exporterBuilder tracer.TracerExporterBuilder) OpenTelemetryOption {
	return OpenTelemetryOptionFunc(func(o *OpenTelemetry) {
		o.opts.tracerOptions = append(o.opts.tracerOptions, tracer.WithExporterBuilder(exporterBuilder))
	})
}
