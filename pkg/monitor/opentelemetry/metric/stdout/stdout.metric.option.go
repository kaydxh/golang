package stdout

func WithMetricPrettyPrint(prettyPrintrl bool) StdoutExporterBuilderOption {
	return StdoutExporterBuilderOptionFunc(func(m *StdoutExporterBuilder) {
		m.opts.prettyPrint = prettyPrintrl
	})
}
