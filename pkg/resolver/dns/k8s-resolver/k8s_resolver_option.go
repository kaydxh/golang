package k8sdns

// A K8sDNSResolverOption sets options.
type K8sDNSResolverOption interface {
	apply(*K8sDNSResolver)
}

// EmptyK8sDNSResolverUrlOption does not alter the K8sDNSResolveruration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyK8sDNSResolverOption struct{}

func (EmptyK8sDNSResolverOption) apply(*K8sDNSResolver) {}

// K8sDNSResolverOptionFunc wraps a function that modifies K8sDNSResolver into an
// implementation of the K8sDNSResolverOption interface.
type K8sDNSResolverOptionFunc func(*K8sDNSResolver)

func (f K8sDNSResolverOptionFunc) apply(do *K8sDNSResolver) {
	f(do)
}

// sample code for option, default for nothing to change
func _K8sDNSResolverOptionWithDefault() K8sDNSResolverOption {
	return K8sDNSResolverOptionFunc(func(*K8sDNSResolver) {
		// nothing to change
	})
}
func (o *K8sDNSResolver) ApplyOptions(options ...K8sDNSResolverOption) *K8sDNSResolver {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
