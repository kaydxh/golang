package resolver

// A ResolveOneOption sets options.
type ResolveOneOption interface {
	apply(*ResolveOneOptions)
}

// EmptyResolveOneUrlOption does not alter the ResolveOneuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyResolveOneOption struct{}

func (EmptyResolveOneOption) apply(*ResolveOneOptions) {}

// ResolveOneOptionFunc wraps a function that modifies ResolveOneOptions into an
// implementation of the ResolveOneOption interface.
type ResolveOneOptionFunc func(*ResolveOneOptions)

func (f ResolveOneOptionFunc) apply(do *ResolveOneOptions) {
	f(do)
}

// sample code for option, default for nothing to change
func _ResolveOneOptionWithDefault() ResolveOneOption {
	return ResolveOneOptionFunc(func(*ResolveOneOptions) {
		// nothing to change
	})
}
func (o *ResolveOneOptions) ApplyOptions(options ...ResolveOneOption) *ResolveOneOptions {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
