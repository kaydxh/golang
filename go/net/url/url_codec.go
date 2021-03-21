package url

import "net/url"

type UrlCodec interface {
	Escape(s string) string
	UnEscape(s string) (string, error)
}

//default url codec
type DefaultUrlCodec struct{}

func (u DefaultUrlCodec) Escape(s string) string {
	return url.QueryEscape(s)
}

func (u DefaultUrlCodec) UnEscape(s string) (string, error) {
	return url.QueryUnescape(s)
}
