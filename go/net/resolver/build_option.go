package resolver

// A ResolverBuildOption sets options.
type ResolverBuildOption interface {
	apply(*ResolverBuildOptions)
}

// EmptyResolverBuildUrlOption does not alter the ResolverBuilduration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyResolverBuildOption struct{}

func (EmptyResolverBuildOption) apply(*ResolverBuildOptions) {}

// ResolverBuildOptionFunc wraps a function that modifies ResolverBuildOptions into an
// implementation of the ResolverBuildOption interface.
type ResolverBuildOptionFunc func(*ResolverBuildOptions)

func (f ResolverBuildOptionFunc) apply(do *ResolverBuildOptions) {
	f(do)
}

// sample code for option, default for nothing to change
func _ResolverBuildOptionWithDefault() ResolverBuildOption {
	return ResolverBuildOptionFunc(func(*ResolverBuildOptions) {
		// nothing to change
	})
}
func (o *ResolverBuildOptions) ApplyOptions(options ...ResolverBuildOption) *ResolverBuildOptions {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
