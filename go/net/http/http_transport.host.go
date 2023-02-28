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
	"net/http"

	resolve_ "github.com/kaydxh/golang/go/net/resolver/resolve"
)

func RequestWithContextTargetHost(req *http.Request, target *Host) *http.Request {
	if target == nil {
		return req
	}
	return req.WithContext(WithContextHost(req.Context(), target))
}

func TargetHostFuncFromContext(req *http.Request) error {
	host := FromContextHost(req.Context())
	if host == nil || host.HostTarget == "" {
		return nil
	}
	if req.URL == nil {
		return nil
	}

	if host.HostTarget == "" {
		return nil
	}

	// replace host of host if target of host if resolved
	address, err := resolve_.ResolveOne(req.Context(), host.HostTarget)
	if err != nil {
		return err
	}
	if address.Addr != "" {
		req.URL.Host = address.Addr
	}
	host.HostTargetAddrResolved = address
	if host.ReplaceHostInRequest {
		req.Host = req.URL.Host
	}
	return nil
}

func RoundTripperWithTarget(rt http.RoundTripper) http.RoundTripper {
	return RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
		err = TargetHostFuncFromContext(req)
		if err != nil {
			return nil, err
		}
		return rt.RoundTrip(req)
	})
}

var DefaultTransportInsecureWithHost = RoundTripperWithTarget(DefaultTransportInsecure)
