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

/*
func WithProxyTargetAddr(addr string) ClientOption {
	proxyFunc := func(req *http.Request) (*url.URL, error) {
		proxyURL, err := url.Parse("http://" + addr)
		if err != nil {
			return proxyURL, nil
		}

		if addr != "" {
			proxyURL.Host = addr
		}

		//	http.ProxyURL(proxyURL)

		return proxyURL, nil
	}

	return ClientOptionFunc(func(c *Client) {
		c.opts.proxy = proxyFunc
	})
}
*/

// http://xxx:yyy@goproxy.com
func WithProxy(proxyURL string) ClientOption {
	return ClientOptionFunc(func(c *Client) {
		c.opts.proxy = proxyURL
	})
}

//dns:///ai-media-1256936300.cos.ap-guangzhou.myqcloud.com
func WithProxyTarget(target string) ClientOption {
	return ClientOptionFunc(func(c *Client) {
		c.opts.proxyTarget = target
	})
}

// WithLogger
func WithLogger(l *log.Logger) ClientOption {
	return ClientOptionFunc(func(c *Client) {
		c.opts.ErrorLog = l
	})
}
