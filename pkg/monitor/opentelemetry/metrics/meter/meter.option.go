package meter

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
