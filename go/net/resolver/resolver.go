package resolver

import (
	"fmt"

	"context"
)

// ResolveNowOptions includes additional information for ResolveNow.
type ResolveNowOptions struct{}

// Resolver watches for the updates on the specified target.
// Updates include address updates and service config updates.
type Resolver interface {
	// It could be called multiple times concurrently.
	ResolveOne(opts ...ResolveOneOption) (Address, error)

	ResolveAll(opts ...ResolveAllOption) ([]Address, error)

	// ResolveNow will be called by gRPC to try to resolve the target name
	// again. It's just a hint, resolver can ignore this if it's not necessary.
	//
	// It could be called multiple times concurrently.
	ResolveNow(opts ...ResolveNowOption)
	// Close closes the resolver.
	Close()
}

type Resolver_PickMode int32

const (
	Resolver_pick_mode_random Resolver_PickMode = 0
	Resolver_pick_mode_first  Resolver_PickMode = 1
)

type Resolver_IPType int32

const (
	Resolver_ip_type_all Resolver_IPType = 0
	Resolver_ip_type_v4  Resolver_IPType = 1
	Resolver_ip_type_v6  Resolver_IPType = 2
)

type ResolveOneOptions struct {
	PickMode Resolver_PickMode
	IPType   Resolver_IPType
}

type ResolveAllOptions struct {
	IPType Resolver_IPType
}

func GetResolver(ctx context.Context, target string, opts ...ResolverBuildOption) (Resolver, error) {
	var opt ResolverBuildOptions
	opt.ApplyOptions(opts...)
	targetInfo, err := ParseTarget(target)
	if err != nil {
		return nil, fmt.Errorf("target[%v] is invalid", targetInfo.Scheme)
	}
	builder := Get(targetInfo.Scheme)
	if builder == nil {
		return nil, fmt.Errorf("scheme[%v] was not registered", targetInfo.Scheme)
	}

	return builder.Build(targetInfo)
}
