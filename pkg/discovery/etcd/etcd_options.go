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
package etcd

// A EtcdKVOption sets options.
type EtcdKVOption interface {
	apply(*EtcdKV)
}

// EmptyEtcdKVUrlOption does not alter the EtcdKVuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyEtcdKVOption struct{}

func (EmptyEtcdKVOption) apply(*EtcdKV) {}

// EtcdKVOptionFunc wraps a function that modifies EtcdKV into an
// implementation of the EtcdKVOption interface.
type EtcdKVOptionFunc func(*EtcdKV)

func (f EtcdKVOptionFunc) apply(do *EtcdKV) {
	f(do)
}

// sample code for option, default for nothing to change
func _EtcdKVOptionWithDefault() EtcdKVOption {
	return EtcdKVOptionFunc(func(*EtcdKV) {
		// nothing to change
	})
}
func (o *EtcdKV) ApplyOptions(options ...EtcdKVOption) *EtcdKV {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}