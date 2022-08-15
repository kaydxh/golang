package http

import (
	"net/http"

	url_ "github.com/kaydxh/golang/go/net/url"
)

func RequestWithProxyTarget(req *http.Request, target string) error {
	if target == "" {
		return nil
	}

	newUrl, err := url_.ReplaceWithTarget(req.Context(), req.URL, target)
	if err != nil {
		return err
	}
	req.URL = newUrl
	req.Host = newUrl.Host

	return nil
}

func NewClientWithTarget(target string) *Client {
	c, _ := NewClient()
	return c
}
