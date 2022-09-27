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
package logs

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

func WithPrefixName(prefixName string) RotateOption {
	return RotateOptionFunc(func(c *Rotate) {
		c.prefixName = prefixName
	})
}

func WithSuffixName(suffixName string) RotateOption {
	return RotateOptionFunc(func(c *Rotate) {
		c.suffixName = suffixName
	})
}

func WithMaxAge(maxAge time.Duration) RotateOption {
	return RotateOptionFunc(func(c *Rotate) {
		c.maxAge = maxAge
	})
}

func WithMaxCount(maxCount int64) RotateOption {
	return RotateOptionFunc(func(c *Rotate) {
		c.maxCount = maxCount
	})
}

func WithRotateSize(rotateSize int64) RotateOption {
	return RotateOptionFunc(func(c *Rotate) {
		c.rotateSize = rotateSize
	})
}

func WithRotateInterval(rotateInterval time.Duration) RotateOption {
	return RotateOptionFunc(func(c *Rotate) {
		c.rotateInterval = rotateInterval
	})
}

type CallerPrettyfierFunc func(f *runtime.Frame) (function string, file string)

func GenShortCallPrettyfier() CallerPrettyfierFunc {
	return func(f *runtime.Frame) (function string, file string) {
		funcname := path.Base(f.Function)
		dir := path.Dir(f.File)
		filename := filepath.Join(path.Base(dir), path.Base(f.File))
		//filename := path.Base(f.File)
		return fmt.Sprintf("%s()", funcname), fmt.Sprintf("%s:%d", filename, f.Line)
	}
}
