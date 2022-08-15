package url

import (
	"context"
	"net/url"
)

func ReplaceWithTarget(ctx context.Context, u *url.URL, target string) (*url.URL, error) {

	if u == nil {
		return nil, nil
	}
	newUrl := u
	if target == "" {
		return newUrl, nil
	}

	newUrl.Host = target
	return newUrl, nil
}
