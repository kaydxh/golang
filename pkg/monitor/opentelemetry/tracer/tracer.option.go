package tracer

func WithExporterBuilder(builder TracerExporterBuilder) TracerOption {
	return TracerOptionFunc(func(m *Tracer) {
		m.opts.builer = builder
	})
}
