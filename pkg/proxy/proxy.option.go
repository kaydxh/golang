package proxy

func WithRouterPatterns(routerPatterns ...string) ReverseProxyOption {
	return ReverseProxyOptionFunc(func(c *ReverseProxy) {
		c.opts.routerPatterns = append(c.opts.routerPatterns, routerPatterns...)
	})
}
