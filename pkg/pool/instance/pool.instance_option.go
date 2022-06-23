package instance

// A PoolOption sets options.
type PoolOption interface {
	apply(*Pool)
}

// EmptyPoolUrlOption does not alter the Pooluration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyPoolOption struct{}

func (EmptyPoolOption) apply(*Pool) {}

// PoolOptionFunc wraps a function that modifies Pool into an
// implementation of the PoolOption interface.
type PoolOptionFunc func(*Pool)

func (f PoolOptionFunc) apply(do *Pool) {
	f(do)
}

// sample code for option, default for nothing to change
func _PoolOptionWithDefault() PoolOption {
	return PoolOptionFunc(func(*Pool) {
		// nothing to change
	})
}
func (o *Pool) ApplyOptions(options ...PoolOption) *Pool {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
