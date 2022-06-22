package jaeger

// A JaegerExporterBuilderOption sets options.
type JaegerExporterBuilderOption interface {
	apply(*JaegerExporterBuilder)
}

// EmptyJaegerExporterBuilderOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyJaegerExporterBuilderOption struct{}

func (EmptyJaegerExporterBuilderOption) apply(*JaegerExporterBuilder) {}

// JaegerExporterBuilderOptionFunc wraps a function that modifies Client into an
// implementation of the JaegerExporterBuilderOption interface.
type JaegerExporterBuilderOptionFunc func(*JaegerExporterBuilder)

func (f JaegerExporterBuilderOptionFunc) apply(do *JaegerExporterBuilder) {
	f(do)
}

// sample code for option, default for nothing to change
func _JaegerExporterBuilderOptionWithDefault() JaegerExporterBuilderOption {
	return JaegerExporterBuilderOptionFunc(func(*JaegerExporterBuilder) {
		// nothing to change
	})
}
func (o *JaegerExporterBuilder) ApplyOptions(options ...JaegerExporterBuilderOption) *JaegerExporterBuilder {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
