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
package resolve

import (
	"context"

	"github.com/kaydxh/golang/go/net/resolver"
	_ "github.com/kaydxh/golang/go/net/resolver/dns"
	_ "github.com/kaydxh/golang/go/net/resolver/passthrough"
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
