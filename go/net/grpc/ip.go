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
