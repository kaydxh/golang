package metric

import "time"

func WithPushExporter(pushExporterBuilder PushExporterBuilder) MeterOption {
	return MeterOptionFunc(func(m *Meter) {
		m.opts.PushExporterBuilder = pushExporterBuilder
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
