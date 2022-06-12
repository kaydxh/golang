package meter

// A MeterOption sets options.
type MeterOption interface {
	apply(*Meter)
}

// EmptyMeterOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyMeterOption struct{}

func (EmptyMeterOption) apply(*Meter) {}

// MeterOptionFunc wraps a function that modifies Client into an
// implementation of the MeterOption interface.
type MeterOptionFunc func(*Meter)

func (f MeterOptionFunc) apply(do *Meter) {
	f(do)
}

// sample code for option, default for nothing to change
func _MeterOptionWithDefault() MeterOption {
	return MeterOptionFunc(func(*Meter) {
		// nothing to change
	})
}
func (o *Meter) ApplyOptions(options ...MeterOption) *Meter {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
