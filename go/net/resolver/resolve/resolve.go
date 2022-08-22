package resolve

import (
	"context"

	"github.com/kaydxh/golang/go/net/resolver"
	_ "github.com/kaydxh/golang/go/net/resolver/dns"
	_ "github.com/kaydxh/golang/go/net/resolver/unix"
)

func ResolveOne(ctx context.Context, target string, opts ...resolver.ResolveOneOption) (resolver.Address, error) {
	r, err := resolver.GetResolver(ctx, target)
	if err != nil {
		return resolver.Address{}, err
	}
	return r.ResolveOne(opts...)
}

func ResolveAll(ctx context.Context, target string, opts ...resolver.ResolveAllOption) ([]resolver.Address, error) {
	r, err := resolver.GetResolver(ctx, target)
	if err != nil {
		return nil, err
	}
	return r.ResolveAll(opts...)
}
