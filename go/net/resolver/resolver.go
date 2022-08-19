package resolver

import (
	"fmt"

	"golang.org/x/net/context"
)

// ResolveNowOptions includes additional information for ResolveNow.
type ResolveNowOptions struct{}

// Resolver watches for the updates on the specified target.
// Updates include address updates and service config updates.
type Resolver interface {
	// It could be called multiple times concurrently.
	ResolveOne(opts ...ResolveOneOption) (Address, error)

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

type ResolveOneOptions struct {
	PickMode Resolver_PickMode
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

func ResolveOne(ctx context.Context, target string, opts ...ResolveOneOption) (Address, error) {
	r, err := GetResolver(ctx, target)
	if err != nil {
		return Address{}, err
	}
	return r.ResolveOne(opts...)
}
