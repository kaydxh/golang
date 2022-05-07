package dns

import "context"

type DNSResolver interface {
	// host domain
	LookupHostIPv4(ctx context.Context, host string) (addrs []string, err error)
}
