package prometheus

// A PrometheusExporterBuilerOption sets options.
type PrometheusExporterBuilerOption interface {
	apply(*PrometheusExporterBuiler)
}

// EmptyPrometheusExporterBuilerOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyPrometheusExporterBuilerOption struct{}

func (EmptyPrometheusExporterBuilerOption) apply(*PrometheusExporterBuiler) {}

// PrometheusExporterBuilerOptionFunc wraps a function that modifies Client into an
// implementation of the PrometheusExporterBuilerOption interface.
type PrometheusExporterBuilerOptionFunc func(*PrometheusExporterBuiler)

func (f PrometheusExporterBuilerOptionFunc) apply(do *PrometheusExporterBuiler) {
	f(do)
}

// sample code for option, default for nothing to change
func _PrometheusExporterBuilerOptionWithDefault() PrometheusExporterBuilerOption {
	return PrometheusExporterBuilerOptionFunc(func(*PrometheusExporterBuiler) {
		// nothing to change
	})
}
func (o *PrometheusExporterBuiler) ApplyOptions(options ...PrometheusExporterBuilerOption) *PrometheusExporterBuiler {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
