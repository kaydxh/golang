package http

import (
	"log"
	"net"
	"net/http"
	"time"
)

type Client struct {
	http.Client
	opts struct {
		timeout               int
		dialTimeout           int
		responseHeaderTimeout int
		maxIdleConns          int
		idleConnTimeout       int
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
		c.Client.Timeout = time.Duration(c.opts.timeout) * time.Second
	}
	if c.opts.dialTimeout != 0 {
		transport.Dial = func(network, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(
				network,
				addr,
				time.Duration(c.opts.dialTimeout)*time.Second,
			)
			if nil != err {
				return nil, err
			}
			return conn, nil
		}
	}

	if c.opts.responseHeaderTimeout != 0 {
		transport.ResponseHeaderTimeout = time.Duration(
			c.opts.responseHeaderTimeout,
		) * time.Second
	}
	if c.opts.maxIdleConns != 0 {
		transport.MaxIdleConns = c.opts.maxIdleConns
	}
	if c.opts.idleConnTimeout != 0 {
		transport.IdleConnTimeout = time.Duration(
			c.opts.idleConnTimeout,
		) * time.Second
	}
	if c.opts.disableKeepAlives {
		transport.DisableKeepAlives = c.opts.disableKeepAlives
	}
	c.Transport = transport

	/*
		cli := &http.Client{
			Transport: transport,
		}
	*/

	return c, nil
}

func (c *Client) Get(url string) (resp *http.Response, err error) {
	return c.Client.Get(url)
}

func (c *Client) logf(format string, args ...interface{}) {
	if c.opts.ErrorLog != nil {
		c.opts.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}
