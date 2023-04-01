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
package rotatefile

import (
	"time"
)

func WithPrefixName(prefixName string) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.prefixName = prefixName
	})
}

func WithFileTimeLayout(fileTimeLayout string) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.fileTimeLayout = fileTimeLayout
	})
}

func WithSuffixName(subfixName string) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.subfixName = subfixName
	})
}

func WithMaxAge(maxAge time.Duration) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.maxAge = maxAge
	})
}

func WithMaxCount(maxCount int64) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.maxCount = maxCount
	})
}

func WithRotateSize(rotateSize int64) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.rotateSize = rotateSize
	})
}

func WithRotateInterval(rotateInterval time.Duration) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.rotateInterval = rotateInterval
	})
}

func WithRotateCallback(callback EventCallbackFunc) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.rotateCallbackFunc = callback
	})
}
