package stdout

import "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

func WithPrettyPrint(prettyPrint bool) StdoutExporterBuilderOption {
	return StdoutExporterBuilderOptionFunc(func(m *StdoutExporterBuilder) {
		if prettyPrint {
			m.opts.stdouttraceOpts = append(m.opts.stdouttraceOpts, stdouttrace.WithPrettyPrint())
		}
	})
}
