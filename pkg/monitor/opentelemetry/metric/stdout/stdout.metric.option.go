package stdout

import "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"

func WithPrettyPrint(prettyPrint bool) StdoutExporterBuilderOption {
	return StdoutExporterBuilderOptionFunc(func(m *StdoutExporterBuilder) {
		if prettyPrint {
			m.opts.stdoutmetricOpts = append(m.opts.stdoutmetricOpts, stdoutmetric.WithPrettyPrint())
		}
	})
}