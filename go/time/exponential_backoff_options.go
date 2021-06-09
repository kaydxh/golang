package time

// A ExponentialBackOffOption sets options.
type ExponentialBackOffOption interface {
	apply(*ExponentialBackOff)
}

// EmptyExponentialBackOffUrlOption does not alter the ExponentialBackOffuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyExponentialBackOffOption struct{}

func (EmptyExponentialBackOffOption) apply(*ExponentialBackOff) {}

// ExponentialBackOffOptionFunc wraps a function that modifies ExponentialBackOff into an
// implementation of the ExponentialBackOffOption interface.
type ExponentialBackOffOptionFunc func(*ExponentialBackOff)

func (f ExponentialBackOffOptionFunc) apply(do *ExponentialBackOff) {
	f(do)
}

// sample code for option, default for nothing to change
func _ExponentialBackOffOptionWithDefault() ExponentialBackOffOption {
	return ExponentialBackOffOptionFunc(func(*ExponentialBackOff) {
		// nothing to change
	})
}
func (o *ExponentialBackOff) ApplyOptions(options ...ExponentialBackOffOption) *ExponentialBackOff {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
