package resolver

import (
	"context"
	"net"
)

// Address represents a server the client connects to.
//
// Experimental
//
// Notice: This type is EXPERIMENTAL and may be changed or removed in a
// later release.
type Address struct {
	// Addr is the server address on which a connection will be established.
	Addr string

	// ServerName is the name of this address.
	// If non-empty, the ServerName is used as the transport certification authority for
	// the address, instead of the hostname from the Dial target string. In most cases,
	// this should not be set.
	//
	// If Type is GRPCLB, ServerName should be the name of the remote load
	// balancer, not the name of the backend.
	//
	// WARNING: ServerName must only be populated with trusted values. It
	// is insecure to populate it with data from untrusted inputs since untrusted
	// values could be used to bypass the authority checks performed by TLS.
	ServerName string
}

// State contains the current Resolver state relevant to the ClientConn.
type State struct {
	// Addresses is the latest set of resolved addresses for the target.
	Addresses []Address
}

// ClientConn contains the callbacks for resolver to notify any updates
// to the gRPC ClientConn.
//
// This interface is to be implemented by gRPC. Users should not need a
// brand new implementation of this interface. For the situations like
// testing, the new implementation should embed this interface. This allows
// gRPC to add new methods to this interface.
type ClientConn interface {
	// UpdateState updates the state of the ClientConn appropriately.
	UpdateState(State) error
	// ReportError notifies the ClientConn that the Resolver encountered an
	// error.  The ClientConn will notify the load balancer and begin calling
	// ResolveNow on the Resolver with exponential backoff.
	ReportError(error)
}

type ResolverBuildOptions struct {
	// if cc is not nil, will update client connection state when resolved
	Cc ClientConn
	// Dialer is the custom dialer used by the ClientConn for dialling the
	// target gRPC service (set via WithDialer). In cases where a name
	// resolution service requires the same dialer, the resolver may use this
	// field. In most cases though, it is not appropriate, and this field may
	// be ignored.
	Dialer func(context.Context, string) (net.Conn, error)

	SkipUnixColonParsing bool
}

/*
type ResolverBuild struct {
	Opts ResolverBuildOptions
}

func NewResolverBuild() {
}
*/

// Builder creates a resolver that will be used to watch name resolution updates.
type Builder interface {
	// Build creates a new resolver for the given target.
	//
	// gRPC dial calls Build synchronously, and fails if the returned error is
	// not nil.
	Build(target Target, opts ...ResolverBuildOption) (Resolver, error)
	// Scheme returns the scheme supported by this resolver.
	// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
	Scheme() string
}
