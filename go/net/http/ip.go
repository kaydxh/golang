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
