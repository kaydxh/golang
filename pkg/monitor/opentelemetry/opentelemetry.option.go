package opentelemetry

import "github.com/kaydxh/golang/pkg/monitor/opentelemetry/metrics/meter"

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
