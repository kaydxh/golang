package opentelemetry

import (
	"time"

	"github.com/kaydxh/golang/pkg/monitor/opentelemetry/metric"
	"github.com/kaydxh/golang/pkg/monitor/opentelemetry/tracer"
)

func WithMeterExporter(exporterBuilder metric.ExporterBuilder) OpenTelemetryOption {
	return OpenTelemetryOptionFunc(func(o *OpenTelemetry) {
		o.opts.meterOptions = append(o.opts.meterOptions, metric.WithExporter(exporterBuilder))

	})
}

func WithMeterPullExporter(pullExporterBuilder metric.PullExporterBuilder) OpenTelemetryOption {
	return OpenTelemetryOptionFunc(func(o *OpenTelemetry) {
		o.opts.meterOptions = append(o.opts.meterOptions, metric.WithPullExporter(pullExporterBuilder))
	})
}

func WithMetricCollectDuration(period time.Duration) OpenTelemetryOption {
	return OpenTelemetryOptionFunc(func(o *OpenTelemetry) {
		o.opts.meterOptions = append(o.opts.meterOptions, metric.WithCollectPeriod(period))
	})
}

func WithTracerExporter(exporterBuilder tracer.TracerExporterBuilder) OpenTelemetryOption {
	return OpenTelemetryOptionFunc(func(o *OpenTelemetry) {
		o.opts.tracerOptions = append(o.opts.tracerOptions, tracer.WithExporterBuilder(exporterBuilder))
	})
}
