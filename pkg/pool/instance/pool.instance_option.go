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
package instance

// A PoolOption sets options.
type PoolOption interface {
	apply(*Pool)
}

// EmptyPoolUrlOption does not alter the Pooluration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyPoolOption struct{}

func (EmptyPoolOption) apply(*Pool) {}

// PoolOptionFunc wraps a function that modifies Pool into an
// implementation of the PoolOption interface.
type PoolOptionFunc func(*Pool)

func (f PoolOptionFunc) apply(do *Pool) {
	f(do)
}

// sample code for option, default for nothing to change
func _PoolOptionWithDefault() PoolOption {
	return PoolOptionFunc(func(*Pool) {
		// nothing to change
	})
}
func (o *Pool) ApplyOptions(options ...PoolOption) *Pool {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
