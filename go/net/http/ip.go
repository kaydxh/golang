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
package http

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// returns IP address from request.
// It will lookup IP in  X-Forwarded-For and X-Real-IP headers.
func GetIPFromRequest(r *http.Request) (net.IP, error) {
	if r == nil {
		return nil, fmt.Errorf("http request is nil")
	}

	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		parts := strings.SplitN(ip, ",", 2)
		part := strings.TrimSpace(parts[0])
		return net.ParseIP(part), nil
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if ip != "" {
		return net.ParseIP(ip), nil
	}

	remoteAddr := strings.TrimSpace(r.RemoteAddr)
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return net.ParseIP(remoteAddr), err
	}

	return net.ParseIP(host), nil
}
