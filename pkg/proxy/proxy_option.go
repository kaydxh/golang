package proxy

// A ProxyOption sets options.
type ProxyOption interface {
	apply(*Proxy)
}

// EmptyProxyOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyProxyOption struct{}

func (EmptyProxyOption) apply(*Proxy) {}

// ProxyOptionFunc wraps a function that modifies Client into an
// implementation of the ProxyOption interface.
type ProxyOptionFunc func(*Proxy)

func (f ProxyOptionFunc) apply(do *Proxy) {
	f(do)
}

// sample code for option, default for nothing to change
func _ProxyOptionWithDefault() ProxyOption {
	return ProxyOptionFunc(func(*Proxy) {
		// nothing to change
	})
}
func (o *Proxy) ApplyOptions(options ...ProxyOption) *Proxy {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
