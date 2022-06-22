package prometheus

func WithMetricUrlPath(url string) PrometheusExporterBuilderOption {
	return PrometheusExporterBuilderOptionFunc(func(m *PrometheusExporterBuilder) {
		m.opts.Url = url
	})
}
