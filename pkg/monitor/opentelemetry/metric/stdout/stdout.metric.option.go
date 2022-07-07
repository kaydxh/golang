package stdout

import "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"

func WithMetricPrettyPrint(prettyPrint bool) StdoutExporterBuilderOption {
	return StdoutExporterBuilderOptionFunc(func(m *StdoutExporterBuilder) {
		if prettyPrint {
			m.opts.stdoutmetricOpts = append(m.opts.stdoutmetricOpts, stdoutmetric.WithPrettyPrint())
		}
	})
}
