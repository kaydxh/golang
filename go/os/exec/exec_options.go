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
package exec

var _default_CommandBuilder_value = func() (val CommandBuilder) { return }()

type CommandBuilderOption interface {
	apply(*CommandBuilder)
}

// EmptyCommandBuilderOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyCommandBuilderOption struct{}

func (EmptyCommandBuilderOption) apply(*CommandBuilder) {}

// CopyDirOptionFunc wraps a function that modifies CopyDir into an
// implementation of the CopyDirOption interface.
type CommandBuilderOptionFunc func(*CommandBuilder)

func (f CommandBuilderOptionFunc) apply(do *CommandBuilder) {
	f(do)
}

// sample code for option, default for nothing to change
func _CommandBuilderOptionWithDefault() CommandBuilderOption {
	return CommandBuilderOptionFunc(func(*CommandBuilder) {
		// nothing to change
	})
}

func (o *CommandBuilder) ApplyOptions(
	options ...CommandBuilderOption,
) *CommandBuilder {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
