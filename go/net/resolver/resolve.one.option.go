package resolver

func WithPickMode(mode Resolver_PickMode) ResolveOneOptionFunc {
	return ResolveOneOptionFunc(func(r *ResolveOneOptions) {
		r.PickMode = mode
	})
}
