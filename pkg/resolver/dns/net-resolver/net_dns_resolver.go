package netdns

import (
	"context"

	net_ "github.com/kaydxh/golang/go/net"
	dns_ "github.com/kaydxh/golang/pkg/resolver/dns"
)

type NetResolver struct {
}

func (n NetResolver) LookupHostIPv4(ctx context.Context, host string) (addrs []string, err error) {
	return net_.LookupHostIPv4WithContext(ctx, host)
}

var DefaultResolver dns_.DNSResolver = NetResolver{}
