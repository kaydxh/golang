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
package http_test

import (
	"testing"
	"time"

	http_ "github.com/kaydxh/golang/go/net/http"
	"golang.org/x/net/context"
	"gotest.tools/v3/assert"
)

func TestHttpClientGet(t *testing.T) {
	ctx := context.Background()
	client, err := http_.NewClient(http_.WithTimeout(5 * time.Second))
	if err != nil {
		t.Fatalf("expect nil, got %v", err)
	}

	testCases := []struct {
		url      string
		expected bool
	}{
		{
			url:      "http://127.0.0.1",
			expected: true,
		},
		{
			url:      "http://127.0.0.2",
			expected: false,
		},
	}

	for _, test := range testCases {
		data, err := client.Get(ctx, test.url)
		if test.expected {
			assert.NilError(t, err)
		} else {
			assert.Assert(t, err != nil) // NotNil
			t.Logf("got %v", err)
		}

		t.Logf("response data: %v", string(data))
	}
}

func TestHttpClientGetWithProxy(t *testing.T) {
	ctx := context.Background()
	client, err := http_.NewClient(http_.WithTimeout(5*time.Second), http_.WithTargetHost("dns:///ai-media-1256936300.cos.ap-guangzhou.myqcloud.com"))
	if err != nil {
		t.Fatalf("expect nil, got %v", err)
	}

	testCases := []struct {
		url      string
		expected bool
	}{
		{
			//url:      "https://ai-media-1256936300.cos.ap-guangzhou.myqcloud.com/find.sh?q-sign-algorithm=sha1&q-ak=AKIDCDyve81SJuISPkMq0IukLg7tupWyoqCg&q-sign-time=1659955959;8640000001659870000&q-key-time=1659955959;8640000001659870000&q-header-list=&q-url-param-list=&q-signature=5695f17ee30c3cd2d37197c773a445eda8a70c8c",
			url:      "http://127.0.0.1/find.sh?q-sign-algorithm=sha1&q-ak=AKIDCDyve81SJuISPkMq0IukLg7tupWyoqCg&q-sign-time=1659955959;8640000001659870000&q-key-time=1659955959;8640000001659870000&q-header-list=&q-url-param-list=&q-signature=5695f17ee30c3cd2d37197c773a445eda8a70c8c",
			expected: true,
		},
	}

	for _, test := range testCases {
		data, err := client.Get(ctx, test.url)
		if test.expected {
			assert.NilError(t, err)
		} else {
			assert.Assert(t, err != nil) // NotNil
			t.Logf("got %v", err)
		}

		t.Logf("response data size: %v", len(data))
	}
}

func TestHttpClientPost(t *testing.T) {
	ctx := context.Background()
	client, err := http_.NewClient(http_.WithTimeout(5 * time.Second))
	if err != nil {
		t.Fatalf("expect nil, got %v", err)
	}

	testCases := []struct {
		url      string
		data     []byte
		expected bool
	}{
		{
			url:      "http://127.0.0.1",
			data:     []byte("hello world test1"),
			expected: true,
		},
		{
			url:      "http://127.0.0.2",
			data:     []byte("hello world test2"),
			expected: false,
		},
	}

	for _, test := range testCases {
		data, err := client.Post(ctx, test.url, "application/text", nil, test.data)
		if test.expected {
			assert.NilError(t, err)
		} else {
			assert.Assert(t, err != nil) // NotNil
			t.Logf("got %v", err)
		}

		t.Logf("response data: %v", string(data))
	}
}
