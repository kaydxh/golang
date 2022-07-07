package stdout

// A StdoutExporterBuilderOption sets options.
type StdoutExporterBuilderOption interface {
	apply(*StdoutExporterBuilder)
}

// EmptyStdoutExporterBuilderOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyStdoutExporterBuilderOption struct{}

func (EmptyStdoutExporterBuilderOption) apply(*StdoutExporterBuilder) {}

// StdoutExporterBuilderOptionFunc wraps a function that modifies Client into an
// implementation of the StdoutExporterBuilderOption interface.
type StdoutExporterBuilderOptionFunc func(*StdoutExporterBuilder)

func (f StdoutExporterBuilderOptionFunc) apply(do *StdoutExporterBuilder) {
	f(do)
}

// sample code for option, default for nothing to change
func _StdoutExporterBuilderOptionWithDefault() StdoutExporterBuilderOption {
	return StdoutExporterBuilderOptionFunc(func(*StdoutExporterBuilder) {
		// nothing to change
	})
}

func (o *StdoutExporterBuilder) ApplyOptions(
	options ...StdoutExporterBuilderOption,
) *StdoutExporterBuilder {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
