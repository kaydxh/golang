package proxy

// A ReverseProxyOption sets options.
type ReverseProxyOption interface {
	apply(*ReverseProxy)
}

// EmptyReverseProxyOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyReverseProxyOption struct{}

func (EmptyReverseProxyOption) apply(*ReverseProxy) {}

// ReverseProxyOptionFunc wraps a function that modifies Client into an
// implementation of the ReverseProxyOption interface.
type ReverseProxyOptionFunc func(*ReverseProxy)

func (f ReverseProxyOptionFunc) apply(do *ReverseProxy) {
	f(do)
}

// sample code for option, default for nothing to change
func _ReverseProxyOptionWithDefault() ReverseProxyOption {
	return ReverseProxyOptionFunc(func(*ReverseProxy) {
		// nothing to change
	})
}
func (o *ReverseProxy) ApplyOptions(options ...ReverseProxyOption) *ReverseProxy {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
