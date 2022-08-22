package resolver

func WithPickMode(mode Resolver_PickMode) ResolveOneOptionFunc {
	return ResolveOneOptionFunc(func(r *ResolveOneOptions) {
		r.PickMode = mode
	})
}

func WithIPTypeForResolverOne(ipType Resolver_IPType) ResolveOneOptionFunc {
	return ResolveOneOptionFunc(func(r *ResolveOneOptions) {
		r.IPType = ipType
	})
}
