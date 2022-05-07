package resolver

func WithResolverType(resolverType Resolver_ResolverType) ResolverQueryOption {
	return ResolverQueryOptionFunc(func(r *ResolverQuery) {
		r.Opts.ResolverType = resolverType
	})
}

func WithLoadBalanceMode(loadBalanceMode Resolver_LoadBalanceMode) ResolverQueryOption {
	return ResolverQueryOptionFunc(func(r *ResolverQuery) {
		r.Opts.LoadBalanceMode = loadBalanceMode
	})
}
