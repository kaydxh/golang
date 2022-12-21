/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
package dns

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	rand_ "github.com/kaydxh/golang/go/math/rand"
	net_ "github.com/kaydxh/golang/go/net"
	"github.com/kaydxh/golang/go/net/resolver"
	time_ "github.com/kaydxh/golang/go/time"
)

// EnableSRVLookups controls whether the DNS resolver attempts to fetch gRPCLB
// addresses from SRV records.  Must not be changed after init time.
var EnableSRVLookups = false

// Globals to stub out in tests. TODO: Perhaps these two can be combined into a
// single variable for testing the resolver?
var (
	newTimer = time.NewTimer
)

func init() {
	resolver.Register(NewBuilder())
}

const (
	defaultPort       = "443"
	defaultDNSSvrPort = "53"
)

var (
	defaultResolver netResolver = net.DefaultResolver
	// To prevent excessive re-resolution, we enforce a rate limit on DNS
	// resolution requests.
	defaultSyncInterval = 30 * time.Second
)

var customAuthorityDialler = func(authority string) func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		var dialer net.Dialer
		return dialer.DialContext(ctx, network, authority)
	}
}

var customAuthorityResolver = func(authority string) (netResolver, error) {
	host, port, err := net_.ParseTarget(authority, defaultDNSSvrPort)
	if err != nil {
		return nil, err
	}

	authorityWithPort := net.JoinHostPort(host, port)

	return &net.Resolver{
		PreferGo: true,
		Dial:     customAuthorityDialler(authorityWithPort),
	}, nil
}

// NewBuilder creates a dnsBuilder which is used to factory DNS resolvers.
func NewBuilder(opts ...dnsBuilderOption) resolver.Builder {
	b := &dnsBuilder{}
	b.ApplyOptions(opts...)
	if b.opts.syncInterval == 0 {
		b.opts.syncInterval = defaultSyncInterval
	}

	return b
}

type dnsBuilder struct {
	opts struct {
		syncInterval time.Duration
	}
}

// Build creates and starts a DNS resolver that watches the name resolution of the target.
func (b *dnsBuilder) Build(target resolver.Target, opts ...resolver.ResolverBuildOption) (resolver.Resolver, error) {
	var opt resolver.ResolverBuildOptions
	opt.ApplyOptions(opts...)
	host, port, err := net_.ParseTarget(target.Endpoint, defaultPort)
	if err != nil {
		return nil, err
	}
	cc := opt.Cc

	// IP address.
	if ipAddr, ok := formatIP(host); ok {
		addr := []resolver.Address{{Addr: ipAddr + ":" + port}}
		if cc != nil {
			cc.UpdateState(resolver.State{Addresses: addr})
		}
		return deadResolver{
			addrs: addr,
		}, nil
	}

	// DNS address (non-IP).
	ctx, cancel := context.WithCancel(context.Background())
	d := &dnsResolver{
		host:         host,
		port:         port,
		syncInterval: b.opts.syncInterval,
		ctx:          ctx,
		cancel:       cancel,
		cc:           cc,
		rn:           make(chan struct{}, 1),
	}

	if target.Authority == "" {
		d.resolver = defaultResolver
	} else {
		d.resolver, err = customAuthorityResolver(target.Authority)
		if err != nil {
			return nil, err
		}
	}

	d.wg.Add(1)
	go d.watcher()
	return d, nil
}

// Scheme returns the naming scheme of this resolver builder, which is "dns".
func (b *dnsBuilder) Scheme() string {
	return "dns"
}

type netResolver interface {
	LookupHost(ctx context.Context, host string) (addrs []string, err error)
	LookupSRV(ctx context.Context, service, proto, name string) (cname string, addrs []*net.SRV, err error)
	LookupTXT(ctx context.Context, name string) (txts []string, err error)
}

// deadResolver is a resolver that does nothing.
type deadResolver struct {
	addrs []resolver.Address
}

func (d deadResolver) ResolveOne(opts ...resolver.ResolveOneOption) (resolver.Address, error) {
	var opt resolver.ResolveOneOptions
	opt.ApplyOptions(opts...)

	addrs, err := d.ResolveAll(resolver.WithIPTypeForResolveAll(opt.IPType))
	if err != nil {
		return resolver.Address{}, err
	}

	switch opt.PickMode {
	case resolver.Resolver_pick_mode_random:
		return addrs[rand_.Intn(len(addrs))], nil
	case resolver.Resolver_pick_mode_first:
		return addrs[0], nil
	default:
		return addrs[rand_.Intn(len(addrs))], nil

	}
}

func (d deadResolver) ResolveAll(opts ...resolver.ResolveAllOption) ([]resolver.Address, error) {
	var opt resolver.ResolveAllOptions
	opt.ApplyOptions(opts...)
	if len(d.addrs) == 0 {
		return nil, fmt.Errorf("resolve target's addresses are empty")
	}

	var pickAddrs []resolver.Address
	if opt.IPType == resolver.Resolver_ip_type_all {
		pickAddrs = d.addrs
	} else {
		for _, addr := range d.addrs {
			v4 := (opt.IPType == resolver.Resolver_ip_type_v4)
			ip, _, _ := net_.SplitHostIntPort(addr.Addr)
			if net_.IsIPv4String(ip) {
				if v4 {
					pickAddrs = append(pickAddrs, addr)
				}
			} else {
				//v6
				if !v4 {
					pickAddrs = append(pickAddrs, addr)
				}
			}
		}
	}
	if len(pickAddrs) == 0 {
		return nil, fmt.Errorf("resolve target's addresses type[%v] are empty", opt.IPType)
	}
	return pickAddrs, nil
}

func (deadResolver) ResolveNow(opts ...resolver.ResolveNowOption) {}

func (deadResolver) Close() {}

// dnsResolver watches for the name resolution update for a non-IP target.
type dnsResolver struct {
	host         string
	port         string
	resolver     netResolver
	syncInterval time.Duration

	ctx    context.Context
	cancel context.CancelFunc
	cc     resolver.ClientConn
	// rn channel is used by ResolveNow() to force an immediate resolution of the target.
	rn chan struct{}
	// wg is used to enforce Close() to return after the watcher() goroutine has finished.
	// Otherwise, data race will be possible. [Race Example] in dns_resolver_test we
	// replace the real lookup functions with mocked ones to facilitate testing.
	// If Close() doesn't wait for watcher() goroutine finishes, race detector sometimes
	// will warns lookup (READ the lookup function pointers) inside watcher() goroutine
	// has data race with replaceNetFunc (WRITE the lookup function pointers).
	wg sync.WaitGroup
}

func (d *dnsResolver) ResolveOne(opts ...resolver.ResolveOneOption) (resolver.Address, error) {
	var opt resolver.ResolveOneOptions
	opt.ApplyOptions(opts...)

	addrs, err := d.ResolveAll(resolver.WithIPTypeForResolveAll(opt.IPType))
	if err != nil {
		return resolver.Address{}, err
	}

	switch opt.PickMode {
	case resolver.Resolver_pick_mode_random:
		return addrs[rand_.Intn(len(addrs))], nil
	case resolver.Resolver_pick_mode_first:
		return addrs[0], nil
	default:
		return addrs[rand_.Intn(len(addrs))], nil

	}
}

func (d *dnsResolver) ResolveAll(opts ...resolver.ResolveAllOption) ([]resolver.Address, error) {
	var opt resolver.ResolveAllOptions
	opt.ApplyOptions(opts...)
	d.ResolveNow()
	addrs, err := d.lookupHost()
	if err != nil {
		return nil, err
	}
	if len(addrs) == 0 {
		return nil, fmt.Errorf("resolve target's addresses are empty")
	}

	var pickAddrs []resolver.Address
	if opt.IPType == resolver.Resolver_ip_type_all {
		pickAddrs = addrs
	} else {
		for _, addr := range addrs {
			v4 := (opt.IPType == resolver.Resolver_ip_type_v4)
			ip, _, _ := net_.SplitHostIntPort(addr.Addr)
			if net_.IsIPv4String(ip) {
				if v4 {
					pickAddrs = append(pickAddrs, addr)
				}
			} else {
				//v6
				if !v4 {
					pickAddrs = append(pickAddrs, addr)
				}
			}
		}
	}
	if len(pickAddrs) == 0 {
		return nil, fmt.Errorf("resolve target's addresses type[%v] are empty", opt.IPType)
	}
	return pickAddrs, nil
}

// ResolveNow invoke an immediate resolution of the target that this dnsResolver watches.
func (d *dnsResolver) ResolveNow(opts ...resolver.ResolveNowOption) {
	select {
	case d.rn <- struct{}{}:
	default:
	}
}

// Close closes the dnsResolver.
func (d *dnsResolver) Close() {
	d.cancel()
	d.wg.Wait()
}

func (d *dnsResolver) watcher() {
	defer d.wg.Done()

	backoff := time_.NewExponentialBackOff()
	for {
		addrs, err := d.lookupHost()
		if d.cc != nil {
			if err != nil {
				// Report error to the underlying grpc.ClientConn.
				d.cc.ReportError(err)
			} else {
				err = d.cc.UpdateState(resolver.State{Addresses: addrs})
			}
		}

		var timer *time.Timer
		if err == nil {
			// Success resolving, wait for the next ResolveNow. However, also wait 30 seconds at the very least
			// to prevent constantly re-resolving.
			backoff.Reset()
			timer = newTimer(d.syncInterval)
			select {
			case <-d.ctx.Done():
				timer.Stop()
				return
			case <-d.rn:
			}
		} else {
			// Poll on an error found in DNS Resolver or an error received from ClientConn.
			actualInterval, _ := backoff.NextBackOff()
			timer = newTimer(actualInterval)
		}
		select {
		case <-d.ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}

func (d *dnsResolver) lookupSRV() ([]resolver.Address, error) {
	if !EnableSRVLookups {
		return nil, nil
	}
	var newAddrs []resolver.Address
	_, srvs, err := d.resolver.LookupSRV(d.ctx, "grpclb", "tcp", d.host)
	if err != nil {
		err = handleDNSError(err, "SRV") // may become nil
		return nil, err
	}
	for _, s := range srvs {
		lbAddrs, err := d.resolver.LookupHost(d.ctx, s.Target)
		if err != nil {
			err = handleDNSError(err, "A") // may become nil
			if err == nil {
				// If there are other SRV records, look them up and ignore this
				// one that does not exist.
				continue
			}
			return nil, err
		}
		for _, a := range lbAddrs {
			ip, ok := formatIP(a)
			if !ok {
				return nil, fmt.Errorf("dns: error parsing A record IP address %v", a)
			}
			addr := ip + ":" + strconv.Itoa(int(s.Port))
			newAddrs = append(newAddrs, resolver.Address{Addr: addr, ServerName: s.Target})
		}
	}
	return newAddrs, nil
}

func handleDNSError(err error, lookupType string) error {
	if dnsErr, ok := err.(*net.DNSError); ok && !dnsErr.IsTimeout && !dnsErr.IsTemporary {
		// Timeouts and temporary errors should be communicated to gRPC to
		// attempt another DNS query (with backoff).  Other errors should be
		// suppressed (they may represent the absence of a TXT record).
		return nil
	}
	if err != nil {
		err = fmt.Errorf("dns: %v record lookup error: %v", lookupType, err)
	}
	return err
}

func (d *dnsResolver) lookupHost() ([]resolver.Address, error) {
	addrs, err := d.resolver.LookupHost(d.ctx, d.host)
	if err != nil {
		err = handleDNSError(err, "A")
		return nil, err
	}
	newAddrs := make([]resolver.Address, 0, len(addrs))
	for _, a := range addrs {
		ip, ok := formatIP(a)
		if !ok {
			return nil, fmt.Errorf("dns: error parsing A record IP address %v", a)
		}
		addr := ip + ":" + d.port
		newAddrs = append(newAddrs, resolver.Address{Addr: addr})
	}
	return newAddrs, nil
}

// formatIP returns ok = false if addr is not a valid textual representation of an IP address.
// If addr is an IPv4 address, return the addr and ok = true.
// If addr is an IPv6 address, return the addr enclosed in square brackets and ok = true.
func formatIP(addr string) (addrIP string, ok bool) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return "", false
	}
	if ip.To4() != nil {
		return addr, true
	}
	return "[" + addr + "]", true
}
