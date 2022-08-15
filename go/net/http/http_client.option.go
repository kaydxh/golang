package http

import (
	"log"
	"net/http"
	"net/url"
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

func WithProxyTarget(addr string) ClientOption {
	return ClientOptionFunc(func(c *Client) {
		c.opts.proxyTarget = addr
	})
}

// WithLogger
func WithLogger(l *log.Logger) ClientOption {
	return ClientOptionFunc(func(c *Client) {
		c.opts.ErrorLog = l
	})
}
