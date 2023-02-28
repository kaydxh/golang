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
package http

import (
	"log"
	"time"
)

func WithTimeout(timeout time.Duration) ClientOption {
	return ClientOptionFunc(func(c *Client) {
		c.opts.timeout = timeout
	})
}

func WithResonseHeaderTimeout(responseHeaderTimeout time.Duration) ClientOption {
	// https://cos.ap-beijing.myqcloud.com
	return ClientOptionFunc(func(c *Client) {
		c.opts.responseHeaderTimeout = responseHeaderTimeout
	})
}

func WithMaxIdleConns(maxIdleConns int) ClientOption {
	return ClientOptionFunc(func(c *Client) {
		c.opts.maxIdleConns = maxIdleConns
	})
}

func WithIdleConnTimeout(idleConnTimeout time.Duration) ClientOption {
	return ClientOptionFunc(func(c *Client) {
		c.opts.idleConnTimeout = idleConnTimeout
	})
}

func WithDisableKeepAlives(disableKeepAlives bool) ClientOption {
	return ClientOptionFunc(func(c *Client) {
		c.opts.disableKeepAlives = disableKeepAlives
	})
}

// http://xxx:yyy@goproxy.com
func WithProxyURL(proxyURL string) ClientOption {
	return ClientOptionFunc(func(c *Client) {
		c.opts.proxyURL = proxyURL
	})
}

//dns:///ai-media-1256936300.cos.ap-guangzhou.myqcloud.com
// redirect
func WithTargetHost(target string) ClientOption {
	return ClientOptionFunc(func(c *Client) {
		c.opts.targetHost = target
	})
}

// replace host in proxy url
func WithProxyHost(target string) ClientOption {
	return ClientOptionFunc(func(c *Client) {
		c.opts.proxyHost = target
	})
}

// WithLogger
func WithLogger(l *log.Logger) ClientOption {
	return ClientOptionFunc(func(c *Client) {
		c.opts.ErrorLog = l
	})
}
