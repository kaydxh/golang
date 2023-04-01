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

func WithMaxFlushBatchSize(batch int) BinlogServiceOption {
	return BinlogServiceOptionFunc(func(c *BinlogService) {
		c.opts.maxFlushBatchSize = batch
	})
}

func WithFlushInterval(flushInterval time.Duration) BinlogServiceOption {
	return BinlogServiceOptionFunc(func(c *BinlogService) {
		c.opts.flushInterval = flushInterval
	})
}

func WithMaxRotateInterval(maxRotateInterval time.Duration) BinlogServiceOption {
	return BinlogServiceOptionFunc(func(c *BinlogService) {
		c.opts.maxRotateInterval = maxRotateInterval
	})
}

func WithMaxRotateSize(maxRotateSize int64) BinlogServiceOption {
	return BinlogServiceOptionFunc(func(c *BinlogService) {
		c.opts.maxRotateSize = maxRotateSize
	})
}
