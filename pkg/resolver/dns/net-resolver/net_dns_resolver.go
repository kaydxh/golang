package netdns

import (
	"net"

	dns_ "github.com/kaydxh/golang/pkg/resolver/dns"
)

var defaultResolver dns_.DNSResolver = net.DefaultResolver
