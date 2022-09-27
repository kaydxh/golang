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
package redis

import "time"

func WithPoolSize(poolSize int) RedisOption {
	return RedisOptionFunc(func(c *RedisClient) {
		c.opts.poolSize = poolSize
	})
}

func WithMinIdleConnections(minIdleConns int) RedisOption {
	return RedisOptionFunc(func(c *RedisClient) {
		c.opts.minIdleConns = minIdleConns
	})
}

func WithDialTimeout(dialTimeout time.Duration) RedisOption {
	return RedisOptionFunc(func(c *RedisClient) {
		c.opts.dialTimeout = dialTimeout
	})
}

func WithReadTimeout(readTimeout time.Duration) RedisOption {
	return RedisOptionFunc(func(c *RedisClient) {
		c.opts.readTimeout = readTimeout
	})
}

func WithWriteTimeout(writeTimeout time.Duration) RedisOption {
	return RedisOptionFunc(func(c *RedisClient) {
		c.opts.writeTimeout = writeTimeout
	})
}

func WithMasterName(masterName string) RedisOption {
	return RedisOptionFunc(func(c *RedisClient) {
		c.opts.masterName = masterName
	})
}
