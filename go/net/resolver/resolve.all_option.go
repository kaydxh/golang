package resolver

// A ResolveAllOption sets options.
type ResolveAllOption interface {
	apply(*ResolveAllOptions)
}

// EmptyResolveAllUrlOption does not alter the ResolveAlluration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyResolveAllOption struct{}

func (EmptyResolveAllOption) apply(*ResolveAllOptions) {}

// ResolveAllOptionFunc wraps a function that modifies ResolveAllOptions into an
// implementation of the ResolveAllOption interface.
type ResolveAllOptionFunc func(*ResolveAllOptions)

func (f ResolveAllOptionFunc) apply(do *ResolveAllOptions) {
	f(do)
}

// sample code for option, default for nothing to change
func _ResolveAllOptionWithDefault() ResolveAllOption {
	return ResolveAllOptionFunc(func(*ResolveAllOptions) {
		// nothing to change
	})
}
func (o *ResolveAllOptions) ApplyOptions(options ...ResolveAllOption) *ResolveAllOptions {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
