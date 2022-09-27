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
