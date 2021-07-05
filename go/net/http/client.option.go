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

// WithLogger
func WithLogger(l *log.Logger) ClientOption {
	return ClientOptionFunc(func(c *Client) {
		c.opts.ErrorLog = l
	})
}
