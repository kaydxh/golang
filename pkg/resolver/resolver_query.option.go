package resolver

import dns_ "github.com/kaydxh/golang/pkg/resolver/dns"

func WithResolverType(resolverType Resolver_ResolverType) ResolverQueryOption {
	return ResolverQueryOptionFunc(func(r *ResolverQuery) {
		r.opts.resolverType = resolverType
	})
}

func WithLoadBalanceMode(loadBalanceMode Resolver_LoadBalanceMode) ResolverQueryOption {
	return ResolverQueryOptionFunc(func(r *ResolverQuery) {
		r.opts.loadBalanceMode = loadBalanceMode
	})
}

func WithNodeGroup(nodeGroup string) ResolverQueryOption {
	return ResolverQueryOptionFunc(func(r *ResolverQuery) {
		r.opts.nodeGroup = nodeGroup
	})
}

func WithNodeUnit(nodeUnit string) ResolverQueryOption {
	return ResolverQueryOptionFunc(func(r *ResolverQuery) {
		r.opts.nodeUnit = nodeUnit
	})
}

func WithResolver(resolver dns_.DNSResolver) ResolverQueryOption {
	return ResolverQueryOptionFunc(func(r *ResolverQuery) {
		r.resolver = resolver
	})
}
