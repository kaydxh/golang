package proxy

func WithProxyMatchedFunc(proxyMatchedFunc ProxyMatchedFunc) ProxyOption {
	return ProxyOptionFunc(func(c *Proxy) {
		c.opts.proxyMatchedFunc = proxyMatchedFunc
	})
}

func WithProxyTargetResolverFunc(proxyTargetResolverFunc ProxyTargetResolverFunc) ProxyOption {
	return ProxyOptionFunc(func(c *Proxy) {
		c.opts.proxyTargetResolverFunc = proxyTargetResolverFunc
	})
}

func WithProxyMode(proxyMode ProxyMode) ProxyOption {
	return ProxyOptionFunc(func(c *Proxy) {
		c.opts.proxyMode = proxyMode
	})
}
