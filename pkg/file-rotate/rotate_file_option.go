package rotatefile

// A RotateFilerOption sets options.
type RotateFilerOption interface {
	apply(*RotateFiler)
}

// EmptyRotateFilerOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyRotateFilerOption struct{}

func (EmptyRotateFilerOption) apply(*RotateFiler) {}

// RotateFilerOptionFunc wraps a function that modifies Client into an
// implementation of the RotateFilerOption interface.
type RotateFilerOptionFunc func(*RotateFiler)

func (f RotateFilerOptionFunc) apply(do *RotateFiler) {
	f(do)
}

// sample code for option, default for nothing to change
func _RotateFilerOptionWithDefault() RotateFilerOption {
	return RotateFilerOptionFunc(func(*RotateFiler) {
		// nothing to change
	})
}
func (o *RotateFiler) ApplyOptions(options ...RotateFilerOption) *RotateFiler {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
