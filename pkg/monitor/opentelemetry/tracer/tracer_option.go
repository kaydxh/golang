package tracer

// A TracerOption sets options.
type TracerOption interface {
	apply(*Tracer)
}

// EmptyTracerOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyTracerOption struct{}

func (EmptyTracerOption) apply(*Tracer) {}

// TracerOptionFunc wraps a function that modifies Client into an
// implementation of the TracerOption interface.
type TracerOptionFunc func(*Tracer)

func (f TracerOptionFunc) apply(do *Tracer) {
	f(do)
}

// sample code for option, default for nothing to change
func _TracerOptionWithDefault() TracerOption {
	return TracerOptionFunc(func(*Tracer) {
		// nothing to change
	})
}
func (o *Tracer) ApplyOptions(options ...TracerOption) *Tracer {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
