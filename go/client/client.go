package client

import (
	"log"
)

type Client struct {
	// options
	opts struct {
		Path     string
		ErrorLog *log.Logger
	}
}

func New(options ...ClientOption) *Client {
	c := &Client{}
	c.ApplyOptions(options...)

	return c
}

func (c *Client) logf(format string, args ...interface{}) {
	if c.opts.ErrorLog != nil {
		c.opts.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}
