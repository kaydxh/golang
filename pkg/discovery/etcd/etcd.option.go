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
package etcd

import (
	"time"
)

func WithDialTimeout(dialTimeout time.Duration) EtcdKVOption {
	return EtcdKVOptionFunc(func(c *EtcdKV) {
		c.opts.DialTimeout = dialTimeout
	})
}

func WithMaxCallRecvMsgSize(msgSize int) EtcdKVOption {
	return EtcdKVOptionFunc(func(c *EtcdKV) {
		c.opts.MaxCallRecvMsgSize = msgSize
	})
}

func WithMaxCallSendMsgSize(msgSize int) EtcdKVOption {
	return EtcdKVOptionFunc(func(c *EtcdKV) {
		c.opts.MaxCallSendMsgSize = msgSize
	})
}

func WithAutoSyncInterval(autoSyncInterval time.Duration) EtcdKVOption {
	return EtcdKVOptionFunc(func(c *EtcdKV) {
		c.opts.AutoSyncInterval = autoSyncInterval
	})
}

func WithWatchPaths(paths []string) EtcdKVOption {
	return EtcdKVOptionFunc(func(c *EtcdKV) {
		c.opts.WatchPaths = paths
	})
}

func WithWatchCreateCallbackFunc(f EventCallbackFunc) EtcdKVOption {
	return EtcdKVOptionFunc(func(c *EtcdKV) {
		c.opts.CreateCallbackFunc = f
	})
}

func WithWatchDeleteCallbackFunc(f EventCallbackFunc) EtcdKVOption {
	return EtcdKVOptionFunc(func(c *EtcdKV) {
		c.opts.DeleteCallbackFunc = f
	})
}

func WithLockTTL(ttl time.Duration) EtcdKVOption {
	return EtcdKVOptionFunc(func(c *EtcdKV) {
		c.opts.LockTTL = ttl
	})
}

func WithLockPrefixPath(prefix string) EtcdKVOption {
	return EtcdKVOptionFunc(func(c *EtcdKV) {
		c.opts.LockPrefixPath = prefix
	})
}

func WithLockKey(key string) EtcdKVOption {
	return EtcdKVOptionFunc(func(c *EtcdKV) {
		c.opts.LockKey = key
	})
}
