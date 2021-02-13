package client

// A ClientOption sets options.
type ClientOption interface {
	apply(*Client)
}

// EmptyClientUrlOption does not alter the Clienturation. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyClientOption struct{}

func (EmptyClientOption) apply(*Client) {}

// ClientOptionFunc wraps a function that modifies Client into an
// implementation of the ClientOption interface.
type ClientOptionFunc func(*Client)

func (f ClientOptionFunc) apply(do *Client) {
	f(do)
}

// sample code for option, default for nothing to change
func _ClientOptionWithDefault() ClientOption {
	return ClientOptionFunc(func(*Client) {
		// nothing to change
	})
}
func (o *Client) ApplyOptions(options ...ClientOption) *Client {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
