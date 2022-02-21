package proxy

func WithRouterPatterns(routerPatterns ...string) ProxyOption {
	return ProxyOptionFunc(func(c *Proxy) {
		c.opts.routerPatterns = append(c.opts.routerPatterns, routerPatterns...)
	})
}

func WithTargetUrl(targetUrl string) ProxyOption {
	return ProxyOptionFunc(func(c *Proxy) {
		c.opts.targetUrl = targetUrl
	})
}

func WithMatchRouterFunc(matchRouter MatchRouterFunc) ProxyOption {
	return ProxyOptionFunc(func(c *Proxy) {
		c.opts.matchRouter = matchRouter
	})
}

func WithProxyMode(proxyMode ProxyMode) ProxyOption {
	return ProxyOptionFunc(func(c *Proxy) {
		c.opts.proxyMode = proxyMode
	})
}
