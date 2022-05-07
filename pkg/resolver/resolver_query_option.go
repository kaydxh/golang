package resolver

// A ResolverQueryOption sets options.
type ResolverQueryOption interface {
	apply(*ResolverQuery)
}

// EmptyResolverQueryUrlOption does not alter the ResolverQueryuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyResolverQueryOption struct{}

func (EmptyResolverQueryOption) apply(*ResolverQuery) {}

// ResolverQueryOptionFunc wraps a function that modifies ResolverQuery into an
// implementation of the ResolverQueryOption interface.
type ResolverQueryOptionFunc func(*ResolverQuery)

func (f ResolverQueryOptionFunc) apply(do *ResolverQuery) {
	f(do)
}

// sample code for option, default for nothing to change
func _ResolverQueryOptionWithDefault() ResolverQueryOption {
	return ResolverQueryOptionFunc(func(*ResolverQuery) {
		// nothing to change
	})
}
func (o *ResolverQuery) ApplyOptions(options ...ResolverQueryOption) *ResolverQuery {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
