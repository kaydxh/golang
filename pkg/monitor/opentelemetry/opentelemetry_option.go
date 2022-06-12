package opentelemetry

// A OpenTelemetryOption sets options.
type OpenTelemetryOption interface {
	apply(*OpenTelemetry)
}

// EmptyOpenTelemetryOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyOpenTelemetryOption struct{}

func (EmptyOpenTelemetryOption) apply(*OpenTelemetry) {}

// OpenTelemetryOptionFunc wraps a function that modifies Client into an
// implementation of the OpenTelemetryOption interface.
type OpenTelemetryOptionFunc func(*OpenTelemetry)

func (f OpenTelemetryOptionFunc) apply(do *OpenTelemetry) {
	f(do)
}

// sample code for option, default for nothing to change
func _OpenTelemetryOptionWithDefault() OpenTelemetryOption {
	return OpenTelemetryOptionFunc(func(*OpenTelemetry) {
		// nothing to change
	})
}
func (o *OpenTelemetry) ApplyOptions(options ...OpenTelemetryOption) *OpenTelemetry {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
