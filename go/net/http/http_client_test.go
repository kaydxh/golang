package http_test

import (
	"testing"
	"time"

	http_ "github.com/kaydxh/golang/go/net/http"
	"gotest.tools/v3/assert"
)

func TestHttpClientGet(t *testing.T) {
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
		data, err := client.Get(test.url)
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
	//client, err := http_.NewClient(http_.WithTimeout(5*time.Second), http_.WithProxyTargetAddr("ai-media-1256936300.cos.ap-guangzhou.myqcloud.com"))
	client, err := http_.NewClient(http_.WithTimeout(5*time.Second), http_.WithProxyTarget("dns:///ai-media-1256936300.cos.ap-guangzhou.myqcloud.com"))
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
		data, err := client.Get(test.url)
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
		data, err := client.Post(test.url, "application/text", nil, test.data)
		if test.expected {
			assert.NilError(t, err)
		} else {
			assert.Assert(t, err != nil) // NotNil
			t.Logf("got %v", err)
		}

		t.Logf("response data: %v", string(data))
	}
}
