package url

import (
	"context"
	"net/url"

	resolve_ "github.com/kaydxh/golang/go/net/resolver/resolve"
)

func ResolveWithTarget(ctx context.Context, u *url.URL, target string) (*url.URL, error) {

	if u == nil {
		return nil, nil
	}
	newUrl := u
	if target == "" {
		return newUrl, nil
	}

	addr, err := resolve_.ResolveOne(ctx, target)
	if err != nil {
		return nil, err
	}
	newUrl.Host = addr.Addr
	return newUrl, nil
}
