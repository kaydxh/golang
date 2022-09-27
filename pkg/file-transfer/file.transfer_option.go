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
package filetransfer

// A FileTransferOption sets options.
type FileTransferOption interface {
	apply(*FileTransfer)
}

// EmptyFileTransferUrlOption does not alter the FileTransferuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyFileTransferOption struct{}

func (EmptyFileTransferOption) apply(*FileTransfer) {}

// FileTransferOptionFunc wraps a function that modifies FileTransfer into an
// implementation of the FileTransferOption interface.
type FileTransferOptionFunc func(*FileTransfer)

func (f FileTransferOptionFunc) apply(do *FileTransfer) {
	f(do)
}

// sample code for option, default for nothing to change
func _FileTransferOptionWithDefault() FileTransferOption {
	return FileTransferOptionFunc(func(*FileTransfer) {
		// nothing to change
	})
}
func (o *FileTransfer) ApplyOptions(options ...FileTransferOption) *FileTransfer {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
