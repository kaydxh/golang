package dns

import "context"

type DNSResolver interface {
	// host domain
	LookupHost(ctx context.Context, host string) (addrs []string, err error)
}
