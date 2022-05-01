package pool

// A PoolOptions sets options.
type PoolOptions interface {
	apply(*Pool)
}

// EmptyPoolUrlOption does not alter the Pooluration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyPoolOptions struct{}

func (EmptyPoolOptions) apply(*Pool) {}

// PoolOptionsFunc wraps a function that modifies Pool into an
// implementation of the PoolOptions interface.
type PoolOptionsFunc func(*Pool)

func (f PoolOptionsFunc) apply(do *Pool) {
	f(do)
}

// sample code for option, default for nothing to change
func _PoolOptionsWithDefault() PoolOptions {
	return PoolOptionsFunc(func(*Pool) {
		// nothing to change
	})
}
func (o *Pool) ApplyOptions(options ...PoolOptions) *Pool {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
