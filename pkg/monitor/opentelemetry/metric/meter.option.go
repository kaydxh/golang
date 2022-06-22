package metric

import "time"

func WithExporter(exporterBuilder ExporterBuilder) MeterOption {
	return MeterOptionFunc(func(m *Meter) {
		m.opts.ExporterBuilder = exporterBuilder
	})
}

func WithPullExporter(pullExporterBuilder PullExporterBuilder) MeterOption {
	return MeterOptionFunc(func(m *Meter) {
		m.opts.PullExporterBuilder = pullExporterBuilder
	})
}

func WithCollectPeriod(period time.Duration) MeterOption {
	return MeterOptionFunc(func(m *Meter) {
		m.opts.collectPeriod = period
	})
}
