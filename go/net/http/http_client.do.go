package http

import (
	"net/http"
)

func (c *Client) get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) post(url string, contentType string, headers map[string]string, auth func(r *http.Request) error) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
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

	return c.Do(req)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	err := RequestWithProxyTarget(req, c.opts.proxyTarget)
	if err != nil {
		return nil, err
	}

	return c.Client.Do(req)
}
