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
	}
	return r, nil
}

func (*passthroughBuilder) Scheme() string {
	return scheme
}

type passthroughResolver struct {
	target resolver.Target
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
