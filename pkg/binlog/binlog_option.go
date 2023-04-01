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
package binlog

// A BinlogServiceOption sets options.
type BinlogServiceOption interface {
	apply(*BinlogService)
}

// EmptyBinlogServiceUrlOption does not alter the BinlogServiceuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyBinlogServiceOption struct{}

func (EmptyBinlogServiceOption) apply(*BinlogService) {}

// BinlogServiceOptionFunc wraps a function that modifies BinlogService into an
// implementation of the BinlogServiceOption interface.
type BinlogServiceOptionFunc func(*BinlogService)

func (f BinlogServiceOptionFunc) apply(do *BinlogService) {
	f(do)
}

// sample code for option, default for nothing to change
func _BinlogServiceOptionWithDefault() BinlogServiceOption {
	return BinlogServiceOptionFunc(func(*BinlogService) {
		// nothing to change
	})
}
func (o *BinlogService) ApplyOptions(options ...BinlogServiceOption) *BinlogService {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
