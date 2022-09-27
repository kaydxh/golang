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

// A ResolverBuildOption sets options.
type ResolverBuildOption interface {
	apply(*ResolverBuildOptions)
}

// EmptyResolverBuildUrlOption does not alter the ResolverBuilduration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyResolverBuildOption struct{}

func (EmptyResolverBuildOption) apply(*ResolverBuildOptions) {}

// ResolverBuildOptionFunc wraps a function that modifies ResolverBuildOptions into an
// implementation of the ResolverBuildOption interface.
type ResolverBuildOptionFunc func(*ResolverBuildOptions)

func (f ResolverBuildOptionFunc) apply(do *ResolverBuildOptions) {
	f(do)
}

// sample code for option, default for nothing to change
func _ResolverBuildOptionWithDefault() ResolverBuildOption {
	return ResolverBuildOptionFunc(func(*ResolverBuildOptions) {
		// nothing to change
	})
}
func (o *ResolverBuildOptions) ApplyOptions(options ...ResolverBuildOption) *ResolverBuildOptions {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
