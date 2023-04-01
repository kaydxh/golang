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

import "time"

func WitRootPath(rootPath string) BinlogServiceOption {
	return BinlogServiceOptionFunc(func(c *BinlogService) {
		c.opts.rootPath = rootPath
	})
}

func WithPrefixName(prefixName string) BinlogServiceOption {
	return BinlogServiceOptionFunc(func(c *BinlogService) {
		c.opts.prefixName = prefixName
	})
}

func WithSufffixName(suffixName string) BinlogServiceOption {
	return BinlogServiceOptionFunc(func(c *BinlogService) {
		c.opts.suffixName = suffixName
	})
}

func WithFlushBatchSize(batch int) BinlogServiceOption {
	return BinlogServiceOptionFunc(func(c *BinlogService) {
		c.opts.flushBatchSize = batch
	})
}

func WithFlushInterval(flushInterval time.Duration) BinlogServiceOption {
	return BinlogServiceOptionFunc(func(c *BinlogService) {
		c.opts.flushInterval = flushInterval
	})
}

func WithRotateInterval(rotateInterval time.Duration) BinlogServiceOption {
	return BinlogServiceOptionFunc(func(c *BinlogService) {
		c.opts.rotateInterval = rotateInterval
	})
}

func WithRotateSize(rotateSize int64) BinlogServiceOption {
	return BinlogServiceOptionFunc(func(c *BinlogService) {
		c.opts.rotateSize = rotateSize
	})
}
