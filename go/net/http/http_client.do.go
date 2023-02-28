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
	"context"
	"io"
	"net/http"

	logs_ "github.com/kaydxh/golang/pkg/logs"
)

func (c *Client) get(ctx context.Context, url string) (*http.Response, error) {
	return c.HttpDo(ctx, http.MethodGet, url, "", nil, nil, nil)
}

func (c *Client) post(ctx context.Context, url string, contentType string, headers map[string]string,
	auth func(r *http.Request) error,
	body io.Reader,
) (*http.Response, error) {
	return c.HttpDo(ctx, http.MethodPost, url, contentType, headers, auth, body)
}

func (c *Client) put(ctx context.Context, url string, contentType string, headers map[string]string,
	auth func(r *http.Request) error,
	body io.Reader,
) (*http.Response, error) {
	return c.HttpDo(ctx, http.MethodPut, url, contentType, headers, auth, body)
}

func (c *Client) HttpDo(ctx context.Context, method string, url string, contentType string, headers map[string]string,
	auth func(r *http.Request) error,
	body io.Reader,
) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
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

	return c.Do(ctx, req)
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	/*
		err := RequestWithTargetHost(req, c.opts.targetHost)
		if err != nil {
			return nil, err
		}
	*/

	logger := logs_.GetLogger(ctx)
	logger.WithField("target_addr", req.Host).Infof("http do %v", req.URL.Path)

	return c.Client.Do(req)
}
