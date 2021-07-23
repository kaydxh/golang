package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type Client struct {
	http.Client
	opts struct {
		timeout               time.Duration
		dialTimeout           time.Duration
		responseHeaderTimeout time.Duration
		idleConnTimeout       time.Duration
		maxIdleConns          int
		disableKeepAlives     bool

		ErrorLog *log.Logger
	}
}

func NewClient(options ...ClientOption) (*Client, error) {
	c := &Client{}
	c.ApplyOptions(options...)
	//	var transport *http.Transport = http.DefaultTransport
	transport := &http.Transport{}
	if c.opts.timeout != 0 {
		c.Client.Timeout = c.opts.timeout
	}
	if c.opts.dialTimeout != 0 {
		transport.Dial = func(network, addr string) (net.Conn, error) {
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
	c.Transport = transport

	return c, nil
}

func (c *Client) Get(url string) ([]byte, error) {
	r, err := c.Client.Get(url)
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
	url, contentType string,
	headers map[string]string,
	body []byte,
) ([]byte, error) {
	bodyReader := bytes.NewReader(body)
	return c.PostReader(url, contentType, headers, nil, bodyReader)
}

func (c *Client) PostJson(
	url string,
	headers map[string]string,
	body []byte,
) ([]byte, error) {
	bodyReader := bytes.NewReader(body)
	return c.PostReader(url, "application/json", headers, nil, bodyReader)
}

func (c *Client) PostJsonWithAuthorize(
	url string,
	headers map[string]string,
	auth func(r *http.Request) error,
	body []byte,
) ([]byte, error) {
	bodyReader := bytes.NewReader(body)
	return c.PostReader(url, "application/json", headers, auth, bodyReader)
}

func (c *Client) PostReader(
	url, contentType string,
	headers map[string]string,
	auth func(r *http.Request) error,
	body io.Reader,
) ([]byte, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if auth != nil {
		err = auth(req)
		if err != nil {
			return nil, err
		}
	}

	r, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if r.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("http status code: %v", r.StatusCode)
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
