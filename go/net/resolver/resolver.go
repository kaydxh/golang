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
		// support no scheme, just ip:port format
		if targetInfo.Scheme == "" {
			return GetDefault().Build(Target{
				Endpoint: target,
			})
		}

		return nil, fmt.Errorf("target[%v] is invalid", targetInfo)
	}

	builder := Get(targetInfo.Scheme)
	if builder == nil {
		return nil, fmt.Errorf("scheme[%v] was not registered", targetInfo.Scheme)
	}
	r, err := builder.Build(targetInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to build target[%v], err: %v", targetInfo, err)
	}
	return r, nil
}
