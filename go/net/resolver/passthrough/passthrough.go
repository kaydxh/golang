package passthrough

import (
	"github.com/kaydxh/golang/go/net/resolver"
)

const scheme = "passthrough"

type passthroughBuilder struct{}

func (*passthroughBuilder) Build(target resolver.Target, opts ...resolver.ResolverBuildOption) (resolver.Resolver, error) {
	var opt resolver.ResolverBuildOptions
	opt.ApplyOptions(opts...)
	r := &passthroughResolver{
		target: target,
		cc:     opt.Cc,
	}
	r.start()
	return r, nil
}

func (*passthroughBuilder) Scheme() string {
	return scheme
}

type passthroughResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
}

func (r *passthroughResolver) start() {
	if r.cc != nil {
		r.cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: r.target.Endpoint}}})
	}
}

func (r *passthroughResolver) ResolveOne(opts ...resolver.ResolveOneOption) (resolver.Address, error) {
	return resolver.Address{Addr: r.target.Endpoint}, nil
}

func (r *passthroughResolver) ResolveAll(opts ...resolver.ResolveAllOption) ([]resolver.Address, error) {
	return []resolver.Address{{Addr: r.target.Endpoint}}, nil
}

func (r *passthroughResolver) ResolveNow(opts ...resolver.ResolveNowOption) {}

func (*passthroughResolver) Close() {}

func init() {
	resolver.Register(&passthroughBuilder{})
}
