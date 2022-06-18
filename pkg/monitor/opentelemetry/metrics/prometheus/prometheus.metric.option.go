package prometheus

func WithMetricUrlPath(urlPath string) PrometheusExporterBuilerOption {
	return PrometheusExporterBuilerOptionFunc(func(m *PrometheusExporterBuiler) {
		m.opts.UrlPath = urlPath
	})
}
