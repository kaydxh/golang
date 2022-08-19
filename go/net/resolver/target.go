package resolver

import (
	"net/url"
	"strings"
)

// Target represents a target for gRPC, as specified in:
// https://github.com/grpc/grpc/blob/master/doc/naming.md.
// It is parsed from the target string that gets passed into Dial or DialContext by the user. And
// grpc passes it to the resolver and the balancer.
//
// If the target follows the naming spec, and the parsed scheme is registered with grpc, we will
// parse the target string according to the spec. e.g. "dns://some_authority/foo.bar" will be parsed
// into &Target{Scheme: "dns", Authority: "some_authority", Endpoint: "foo.bar"}
//
// If the target does not contain a scheme, we will apply the default scheme, and set the Target to
// be the full target string. e.g. "foo.bar" will be parsed into
// &Target{Scheme: resolver.GetDefaultScheme(), Endpoint: "foo.bar"}.
//
// If the parsed scheme is not registered (i.e. no corresponding resolver available to resolve the
// endpoint), we set the Scheme to be the default scheme, and set the Endpoint to be the full target
// string. e.g. target string "unknown_scheme://authority/endpoint" will be parsed into
// &Target{Scheme: resolver.GetDefaultScheme(), Endpoint: "unknown_scheme://authority/endpoint"}.
type Target struct {
	Scheme    string
	Authority string
	Endpoint  string
	URL       url.URL
}

// parseTarget uses RFC 3986 semantics to parse the given target into a
// resolver.Target struct containing scheme, authority and endpoint. Query
// params are stripped from the endpoint.
func ParseTarget(target string) (Target, error) {
	u, err := url.Parse(target)
	if err != nil {
		return Target{}, err
	}
	// For targets of the form "[scheme]://[authority]/endpoint, the endpoint
	// value returned from url.Parse() contains a leading "/". Although this is
	// in accordance with RFC 3986, we do not want to break existing resolver
	// implementations which expect the endpoint without the leading "/". So, we
	// end up stripping the leading "/" here. But this will result in an
	// incorrect parsing for something like "unix:///path/to/socket". Since we
	// own the "unix" resolver, we can workaround in the unix resolver by using
	// the `URL` field instead of the `Endpoint` field.
	endpoint := u.Path
	if endpoint == "" {
		endpoint = u.Opaque
	}

	endpoint = strings.TrimPrefix(endpoint, "/")
	return Target{
		Scheme:    u.Scheme,
		Authority: u.Host,
		Endpoint:  endpoint,
		URL:       *u,
	}, nil
}

// ParseTarget splits target into a resolver.Target struct containing scheme,
// authority and endpoint. skipUnixColonParsing indicates that the parse should
// not parse "unix:[path]" cases. This should be true in cases where a custom
// dialer is present, to prevent a behavior change.
//
// If target is not a valid scheme://authority/endpoint as specified in
// https://github.com/grpc/grpc/blob/master/doc/naming.md,
// it returns {Endpoint: target}.
/*
func ParseTarget(target string, skipUnixColonParsing bool) (ret Target) {
	var ok bool
	if strings.HasPrefix(target, "unix-abstract:") {
		if strings.HasPrefix(target, "unix-abstract://") {
			// Maybe, with Authority specified, try to parse it
			var remain string
			ret.Scheme, remain, _ = strings_.Split2(target, "://")
			ret.Authority, ret.Endpoint, ok = strings_.Split2(remain, "/")
			if !ok {
				// No Authority, add the "//" back
				ret.Endpoint = "//" + remain
			} else {
				// Found Authority, add the "/" back
				ret.Endpoint = "/" + ret.Endpoint
			}
		} else {
			// Without Authority specified, split target on ":"
			ret.Scheme, ret.Endpoint, _ = strings_.Split2(target, ":")
		}
		return ret
	}
	ret.Scheme, ret.Endpoint, ok = strings_.Split2(target, "://")
	if !ok {
		if strings.HasPrefix(target, "unix:") && !skipUnixColonParsing {
			// Handle the "unix:[local/path]" and "unix:[/absolute/path]" cases,
			// because splitting on :// only handles the
			// "unix://[/absolute/path]" case. Only handle if the dialer is nil,
			// to avoid a behavior change with custom dialers.
			return Target{Scheme: "unix", Endpoint: target[len("unix:"):]}
		}
		return Target{Endpoint: target}
	}
	ret.Authority, ret.Endpoint, ok = strings_.Split2(ret.Endpoint, "/")
	if !ok {
		return Target{Endpoint: target}
	}
	if ret.Scheme == "unix" {
		// Add the "/" back in the unix case, so the unix resolver receives the
		// actual endpoint in the "unix://[/absolute/path]" case.
		ret.Endpoint = "/" + ret.Endpoint
	}
	return ret
}
*/
