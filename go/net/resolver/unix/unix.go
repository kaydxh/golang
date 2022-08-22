package unix

import (
	"fmt"

	"github.com/kaydxh/golang/go/net/resolver"
)

const unixScheme = "unix"
const unixAbstractScheme = "unix-abstract"

type builder struct {
	scheme string
}

func (b *builder) Build(target resolver.Target, opts ...resolver.ResolverBuildOption) (resolver.Resolver, error) {
	if target.Authority != "" {
		return nil, fmt.Errorf("invalid (non-empty) authority: %v", target.Authority)
	}
	var opt resolver.ResolverBuildOptions
	opt.ApplyOptions(opts...)
	addr := resolver.Address{Addr: target.Endpoint}
	if b.scheme == unixAbstractScheme {
		// prepend "\x00" to address for unix-abstract
		addr.Addr = "\x00" + addr.Addr
	}
	cc := opt.Cc
	if cc != nil {
		cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: addr.Addr}}})
	}
	return &nopResolver{addrs: []resolver.Address{{Addr: addr.Addr}}}, nil
}

func (b *builder) Scheme() string {
	return b.scheme
}

type nopResolver struct {
	addrs []resolver.Address
}

func (*nopResolver) ResolveOne(opts ...resolver.ResolveOneOption) (resolver.Address, error) {
	return resolver.Address{}, nil
}

func (*nopResolver) ResolveAll(opts ...resolver.ResolveAllOption) ([]resolver.Address, error) {
	return nil, nil
}

func (*nopResolver) ResolveNow(opts ...resolver.ResolveNowOption) {}

func (*nopResolver) Close() {}

func init() {
	resolver.Register(&builder{scheme: unixScheme})
	resolver.Register(&builder{scheme: unixAbstractScheme})
}
