package base64

import (
	"encoding/base64"
	"fmt"
)

func EncodeString(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func DecodeString(v string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %v", err)
	}

	return string(data), nil
}

func EncodeURL(v string) string {
	return base64.URLEncoding.EncodeToString([]byte(v))
}

func DecodeURL(v string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(v)
	if err != nil {
		return "", fmt.Errorf("base64 decode url failed: %v", err)
	}

	return string(data), nil
}
