package resolver

func WithIPTypeForResolveAll(ipType Resolver_IPType) ResolveAllOptionFunc {
	return ResolveAllOptionFunc(func(r *ResolveAllOptions) {
		r.IPType = ipType
	})
}
