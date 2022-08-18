package unix

import (
	"fmt"

	"google.golang.org/grpc/internal/transport/networktype"
	"google.golang.org/grpc/resolver"
)

const unixScheme = "unix"
const unixAbstractScheme = "unix-abstract"

type builder struct {
	scheme string
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	if target.Authority != "" {
		return nil, fmt.Errorf("invalid (non-empty) authority: %v", target.Authority)
	}
	addr := resolver.Address{Addr: target.Endpoint}
	if b.scheme == unixAbstractScheme {
		// prepend "\x00" to address for unix-abstract
		addr.Addr = "\x00" + addr.Addr
	}
	cc.UpdateState(resolver.State{Addresses: []resolver.Address{networktype.Set(addr, "unix")}})
	return &nopResolver{}, nil
}

func (b *builder) Scheme() string {
	return b.scheme
}

type nopResolver struct {
}

func (*nopResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (*nopResolver) Close() {}

func init() {
	resolver.Register(&builder{scheme: unixScheme})
	resolver.Register(&builder{scheme: unixAbstractScheme})
}
