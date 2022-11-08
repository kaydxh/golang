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
package s3

// A StorageOption sets options.
type StorageOption interface {
	apply(*Storage)
}

// EmptyStorageUrlOption does not alter the Storageuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyStorageOption struct{}

func (EmptyStorageOption) apply(*Storage) {}

// StorageOptionFunc wraps a function that modifies Storage into an
// implementation of the StorageOption interface.
type StorageOptionFunc func(*Storage)

func (f StorageOptionFunc) apply(do *Storage) {
	f(do)
}

// sample code for option, default for nothing to change
func _StorageOptionWithDefault() StorageOption {
	return StorageOptionFunc(func(*Storage) {
		// nothing to change
	})
}
func (o *Storage) ApplyOptions(options ...StorageOption) *Storage {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
