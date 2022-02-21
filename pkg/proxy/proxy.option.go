package proxy

func WithRouterPatterns(routerPatterns ...string) ReverseProxyOption {
	return ReverseProxyOptionFunc(func(c *ReverseProxy) {
		c.opts.routerPatterns = append(c.opts.routerPatterns, routerPatterns...)
	})
}

func WithTargetUrl(targetUrl string) ReverseProxyOption {
	return ReverseProxyOptionFunc(func(c *ReverseProxy) {
		c.opts.targetUrl = targetUrl
	})
}

func WithMatchRouterFunc(matchRouter MatchRouterFunc) ReverseProxyOption {
	return ReverseProxyOptionFunc(func(c *ReverseProxy) {
		c.opts.matchRouter = matchRouter
	})
}

func WithProxyMode(proxyMode ProxyMode) ReverseProxyOption {
	return ReverseProxyOptionFunc(func(c *ReverseProxy) {
		c.opts.proxyMode = proxyMode
	})
}
