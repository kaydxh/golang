package client

func WithClientOptionPath(path string) ClientOption {
	return ClientOptionFunc(func(opt *Client) {
		opt.opts.Path = path
	})
}
