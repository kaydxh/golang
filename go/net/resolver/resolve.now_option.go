package resolver

// A ResolveNowOption sets options.
type ResolveNowOption interface {
	apply(*ResolveNowOptions)
}

// EmptyResolveNowUrlOption does not alter the ResolveNowuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyResolveNowOption struct{}

func (EmptyResolveNowOption) apply(*ResolveNowOptions) {}

// ResolveNowOptionFunc wraps a function that modifies ResolveNowOptions into an
// implementation of the ResolveNowOption interface.
type ResolveNowOptionFunc func(*ResolveNowOptions)

func (f ResolveNowOptionFunc) apply(do *ResolveNowOptions) {
	f(do)
}

// sample code for option, default for nothing to change
func _ResolveNowOptionWithDefault() ResolveNowOption {
	return ResolveNowOptionFunc(func(*ResolveNowOptions) {
		// nothing to change
	})
}
func (o *ResolveNowOptions) ApplyOptions(options ...ResolveNowOption) *ResolveNowOptions {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
