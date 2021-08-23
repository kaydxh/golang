package base64_test

import (
	"testing"

	base64_ "github.com/kaydxh/golang/go/encoding/base64"
	"gotest.tools/v3/assert"
)

func TestString(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
	}{
		{
			name:     "hello word",
			expected: "",
		},
		{
			name:     "http://12306.com?a=%b",
			expected: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			encoded := base64_.EncodeString(testCase.name)
			t.Logf("base64 encode : %v", encoded)
			decoded, err := base64_.DecodeString(encoded)
			if err != nil {
				t.Fatalf("failed to decode string, err: %v", err)
			}

			assert.Equal(t, testCase.name, decoded)
		})
	}
	/*
		content := "hello word"

		encoded := base64_.EncodeString(content)
		t.Logf("base64 encode : %v", encoded)

		decoded, err := base64_.DecodeString(encoded)
		if err != nil {
			t.Fatalf("failed to decode string, err: %v", err)
		}
		t.Logf("base64 decode : %v", decoded)
	*/
	//assert.Equal(t, content, decoded)
}
