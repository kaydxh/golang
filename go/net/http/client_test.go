package http_test

import (
	"testing"

	http_ "github.com/kaydxh/golang/go/net/http"
	"gotest.tools/v3/assert"
)

func TestHttpClientGet(t *testing.T) {
	client, err := http_.NewClient(http_.WithTimeout(5))
	if err != nil {
		t.Errorf("expect nil, got %v", err)
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

func TestHttpClientPost(t *testing.T) {
	client, err := http_.NewClient(http_.WithTimeout(5))
	if err != nil {
		t.Errorf("expect nil, got %v", err)
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
		data, err := client.Post(test.url, "application/text", test.data)
		if test.expected {
			assert.NilError(t, err)
		} else {
			assert.Assert(t, err != nil) // NotNil
			t.Logf("got %v", err)
		}

		t.Logf("response data: %v", string(data))
	}
}
