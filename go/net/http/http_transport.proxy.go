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
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	resolve_ "github.com/kaydxh/golang/go/net/resolver/resolve"
)

func RequestWithContextProxy(req *http.Request, proxy *Proxy) *http.Request {
	if proxy == nil {
		return req
	}
	return req.WithContext(WithContextProxy(req.Context(), proxy))
}

func ProxyFuncFromContextOrEnvironment(req *http.Request) (*url.URL, error) {
	proxy := FromContextProxy(req.Context())
	if proxy == nil || proxy.ProxyUrl == "" {
		return http.ProxyFromEnvironment(req)
	}

	proxyUrl, err := ParseProxyUrl(proxy.ProxyUrl)
	if err != nil {
		return nil, err
	}
	if proxyUrl == nil {
		return nil, nil
	}

	if proxy.ProxyTarget == "" {
		return proxyUrl, nil
	}

	// replace host of proxy if target of proxy if resolved
	address, err := resolve_.ResolveOne(req.Context(), proxy.ProxyTarget)
	if err != nil {
		return nil, err
	}
	if address.Addr != "" {
		proxyUrl.Host = address.Addr
	}
	proxy.ProxyAddrResolved = address
	return proxyUrl, nil
}

var DefaultTransportInsecureWithProxy http.RoundTripper = &http.Transport{
	Proxy: ProxyFuncFromContextOrEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	// skip verify for https
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
