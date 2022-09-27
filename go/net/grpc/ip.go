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
package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// returns IP address from grpc context.
// It will lookup IP in  X-Forwarded-For and X-Real-IP headers, if both
// get empty, lookup IP from peer context
func GetIPFromContext(ctx context.Context) (net.IP, error) {
	// FromIncomingContext returns the incoming metadata in ctx if it exists.
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		peerAddr := md.Get("x-real-ip")
		if len(peerAddr) > 0 {
			return net.ParseIP(peerAddr[0]), nil
		}

		peerAddr = md.Get("x-forwarded-for")
		if len(peerAddr) > 0 {
			return net.ParseIP(peerAddr[0]), nil
		}
	}

	//if use proxy, only return proxy address
	// FromContext returns the peer information in ctx if it exists.
	peerAddr, ok := peer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unexpected context")
	}

	if peerAddr.Addr == net.Addr(nil) {
		return nil, fmt.Errorf("unexpected err: peer address is nil")
	}

	host, _, err := net.SplitHostPort(peerAddr.Addr.String())
	if err != nil {
		return nil, fmt.Errorf("invalid peer host: %v, err: %v", peerAddr.Addr.String(), err)
	}

	return net.ParseIP(host), nil
}
