package prometheus

// A PrometheusExporterBuilderOption sets options.
type PrometheusExporterBuilderOption interface {
	apply(*PrometheusExporterBuilder)
}

// EmptyPrometheusExporterBuilderOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyPrometheusExporterBuilderOption struct{}

func (EmptyPrometheusExporterBuilderOption) apply(*PrometheusExporterBuilder) {}

// PrometheusExporterBuilderOptionFunc wraps a function that modifies Client into an
// implementation of the PrometheusExporterBuilderOption interface.
type PrometheusExporterBuilderOptionFunc func(*PrometheusExporterBuilder)

func (f PrometheusExporterBuilderOptionFunc) apply(do *PrometheusExporterBuilder) {
	f(do)
}

// sample code for option, default for nothing to change
func _PrometheusExporterBuilderOptionWithDefault() PrometheusExporterBuilderOption {
	return PrometheusExporterBuilderOptionFunc(func(*PrometheusExporterBuilder) {
		// nothing to change
	})
}

func (o *PrometheusExporterBuilder) ApplyOptions(
	options ...PrometheusExporterBuilderOption,
) *PrometheusExporterBuilder {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
