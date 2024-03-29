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
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin/binding"
)

type Client struct {
	http.Client
	opts struct {
		// Timeout specifies a time limit for requests made by this
		// Client. The timeout includes connection time, any
		// redirects, and reading the response body. The timer remains
		// running after Get, Head, Post, or Do return and will
		// interrupt reading of the Response.Body.
		//
		// A Timeout of zero means no timeout.
		//
		// The Client cancels requests to the underlying Transport
		// as if the Request's Context ended.
		//
		// For compatibility, the Client will also use the deprecated
		// CancelRequest method on Transport if found. New
		// RoundTripper implementations should use the Request's Context
		// for cancellation instead of implementing CancelRequest.
		timeout               time.Duration
		dialTimeout           time.Duration
		responseHeaderTimeout time.Duration
		idleConnTimeout       time.Duration
		maxIdleConns          int
		disableKeepAlives     bool

		// Proxy specifies a function to return a proxy for a given
		// Request. If the function returns a non-nil error, the
		// request is aborted with the provided error.
		//
		// The proxy type is determined by the URL scheme. "http",
		// "https", and "socks5" are supported. If the scheme is empty,
		// "http" is assumed.
		//
		// If Proxy is nil or returns a nil *URL, no proxy is used.
		//proxy func(*http.Request) (*url.URL, error)
		// like forward proxy
		proxyURL string

		// proxyHost is host:port, or domain, replace host in proxy
		proxyHost string

		// targetHost is host:port, redirct to it, like reverse proxy
		targetHost string

		ErrorLog *log.Logger
	}
}

func NewClient(options ...ClientOption) (*Client, error) {
	c := &Client{}
	c.ApplyOptions(options...)
	transport := DefaultTransportInsecure
	/*
		transport := &http.Transport{

			// ProxyFromEnvironment returns the URL of the proxy to use for a
			// given request, as indicated by the environment variables
			// HTTP_PROXY, HTTPS_PROXY and NO_PROXY (or the lowercase versions
			// thereof). HTTPS_PROXY takes precedence over HTTP_PROXY for https
			// requests.
			//
			// The environment values may be either a complete URL or a
			// "host[:port]", in which case the "http" scheme is assumed.
			// The schemes "http", "https", and "socks5" are supported.
			// An error is returned if the value is a different form.
			//
			// A nil URL and nil error are returned if no proxy is defined in the
			// environment, or a proxy should not be used for the given request,
			// as defined by NO_PROXY.
			//
			// As a special case, if req.URL.Host is "localhost" (with or without
			// a port number), then a nil URL and nil error will be returned.
			Proxy: http.ProxyFromEnvironment,
			// skip verify for https
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	*/
	if c.opts.timeout != 0 {
		c.Client.Timeout = c.opts.timeout
	}
	if c.opts.dialTimeout != 0 {
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(
				network,
				addr,
				c.opts.dialTimeout,
			)
			if nil != err {
				return nil, err
			}
			return conn, nil
		}
	}

	if c.opts.responseHeaderTimeout != 0 {
		transport.ResponseHeaderTimeout = c.opts.responseHeaderTimeout
	}
	if c.opts.maxIdleConns != 0 {
		transport.MaxIdleConns = c.opts.maxIdleConns
	}
	if c.opts.idleConnTimeout != 0 {
		transport.IdleConnTimeout = c.opts.idleConnTimeout
	}
	if c.opts.disableKeepAlives {
		transport.DisableKeepAlives = c.opts.disableKeepAlives
	}
	c.Transport = RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
		if c.opts.proxyURL != "" || c.opts.proxyHost != "" {
			transport.Proxy = ProxyFuncFromContextOrEnvironment

			proxyUrl := "http://"
			if c.opts.proxyURL != "" {
				proxyUrl = c.opts.proxyURL
			}
			proxy := &Proxy{
				ProxyUrl:    proxyUrl,
				ProxyTarget: c.opts.proxyHost,
			}
			req = RequestWithContextProxy(req, proxy)
		}

		if c.opts.targetHost != "" {
			host := &Host{
				HostTarget:           c.opts.targetHost,
				ReplaceHostInRequest: true,
			}
			req = RequestWithContextTargetHost(req, host)
		}

		return RoundTripperWithTarget(transport).RoundTrip(req)

	})

	return c, nil
}

func (c *Client) Get(ctx context.Context, url string) ([]byte, error) {
	r, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) Post(
	ctx context.Context,
	url, contentType string,
	headers map[string]string,
	body []byte,
) ([]byte, error) {
	bodyReader := bytes.NewReader(body)
	return c.PostReader(ctx, url, contentType, headers, nil, bodyReader)
}

func (c *Client) Put(
	ctx context.Context,
	url, contentType string,
	headers map[string]string,
	body []byte,
) ([]byte, error) {
	bodyReader := bytes.NewReader(body)
	return c.PutReader(ctx, url, contentType, headers, nil, bodyReader)
}

func (c *Client) PostJson(
	ctx context.Context,
	url string,
	headers map[string]string,
	body []byte,
) ([]byte, error) {
	bodyReader := bytes.NewReader(body)
	return c.PostReader(ctx, url, binding.MIMEJSON, headers, nil, bodyReader)
}

func (c *Client) PostPb(
	ctx context.Context,
	url string,
	headers map[string]string,
	body []byte,
) ([]byte, error) {
	bodyReader := bytes.NewReader(body)
	return c.PostReader(ctx, url, binding.MIMEPROTOBUF, headers, nil, bodyReader)
}

func (c *Client) PostJsonWithAuthorize(
	ctx context.Context,
	url string,
	headers map[string]string,
	auth func(r *http.Request) error,
	body []byte,
) ([]byte, error) {
	bodyReader := bytes.NewReader(body)
	return c.PostReader(ctx, url, binding.MIMEJSON, headers, auth, bodyReader)
}

func (c *Client) PostReader(
	ctx context.Context,
	url, contentType string,
	headers map[string]string,
	auth func(r *http.Request) error,
	body io.Reader,
) ([]byte, error) {
	return c.HttpReader(ctx, http.MethodPost, url, contentType, headers, auth, body)
}

func (c *Client) PutReader(
	ctx context.Context,
	url, contentType string,
	headers map[string]string,
	auth func(r *http.Request) error,
	body io.Reader,
) ([]byte, error) {
	return c.HttpReader(ctx, http.MethodPut, url, contentType, headers, auth, body)
}

func (c *Client) HttpReader(
	ctx context.Context,
	method, url, contentType string,
	headers map[string]string,
	auth func(r *http.Request) error,
	body io.Reader,
) ([]byte, error) {
	r, err := c.HttpDo(ctx, method, url, contentType, headers, auth, body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if r.StatusCode >= http.StatusBadRequest {
		return data, fmt.Errorf("http status code: %v", r.StatusCode)
	}

	return data, nil
}

func (c *Client) logf(format string, args ...interface{}) {
	if c.opts.ErrorLog != nil {
		c.opts.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}
