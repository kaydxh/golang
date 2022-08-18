package resolver

import (
	"context"
	"net"
)

func WithBuildDialer(dialer func(context.Context, string) (net.Conn, error)) ResolverBuildOptionFunc {
	return ResolverBuildOptionFunc(func(b *ResolverBuildOptions) {
		b.Dialer = dialer
	})
}
