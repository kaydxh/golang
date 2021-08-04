package logs

// A RotateOption sets options.
type RotateOption interface {
	apply(*Rotate)
}

// EmptyRotateOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyRotateOption struct{}

func (EmptyRotateOption) apply(*Rotate) {}

// RotateOptionFunc wraps a function that modifies Client into an
// implementation of the RotateOption interface.
type RotateOptionFunc func(*Rotate)

func (f RotateOptionFunc) apply(do *Rotate) {
	f(do)
}

// sample code for option, default for nothing to change
func _RotateOptionWithDefault() RotateOption {
	return RotateOptionFunc(func(*Rotate) {
		// nothing to change
	})
}

func (o *Rotate) ApplyOptions(options ...RotateOption) *Rotate {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
