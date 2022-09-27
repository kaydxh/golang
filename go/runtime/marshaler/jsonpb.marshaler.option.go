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
package marshaler

func WithUseProtoNames(useProtoNames bool) JSONPbOption {
	return JSONPbOptionFunc(func(c *JSONPb) {
		c.opts.useProtoNames = useProtoNames
	})
}

func WithUseEnumNumbers(useEnumNumbers bool) JSONPbOption {
	return JSONPbOptionFunc(func(c *JSONPb) {
		c.opts.useEnumNumbers = useEnumNumbers
	})
}

func WithEmitUnpopulated(emitUnpopulated bool) JSONPbOption {
	return JSONPbOptionFunc(func(c *JSONPb) {
		c.opts.emitUnpopulated = emitUnpopulated
	})
}

func WithDiscardUnknown(discardUnknown bool) JSONPbOption {
	return JSONPbOptionFunc(func(c *JSONPb) {
		c.opts.discardUnknown = discardUnknown
	})
}

func WithIndent(indent string) JSONPbOption {
	return JSONPbOptionFunc(func(c *JSONPb) {
		c.opts.indent = indent
	})
}
