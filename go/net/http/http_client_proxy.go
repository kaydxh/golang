package http

import (
	"net/http"

	url_ "github.com/kaydxh/golang/go/net/url"
)

func RequestWithProxyTarget(req *http.Request, target string) error {
	if target == "" {
		return nil
	}

	newUrl, err := url_.ResolveWithTarget(req.Context(), req.URL, target)
	if err != nil {
		return err
	}
	req.URL = newUrl
	req.Host = newUrl.Host

	return nil
}

func NewClientWithProxyTarget(target string, opts ...ClientOption) *Client {
	opts = append(opts, WithProxyTarget(target))
	c, _ := NewClient(opts...)
	return c
}
