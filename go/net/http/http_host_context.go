/*
 *Copyright (c) 2023, kaydxh
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
package http

import (
	"context"
	"fmt"
	"net/url"

	"github.com/kaydxh/golang/go/net/resolver"
)

type hostContextKey struct{}

type Host struct {
	HostTarget           string
	ReplaceHostInRequest bool

	HostTargetAddrResolved resolver.Address
}

func FromContextHost(ctx context.Context) *Host {
	host, _ := ctx.Value(hostContextKey{}).(*Host)
	return host
}

func WithContextHost(ctx context.Context, host *Host) context.Context {
	if host == nil {
		panic("nil host")
	}
	return context.WithValue(ctx, hostContextKey{}, host)
}

func ParseTargetUrl(host string) (*url.URL, error) {
	if host == "" {
		return nil, nil
	}

	hostURL, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("invalid host address %q: %v", host, err)
	}
	return hostURL, nil
}
