package http

// A HandlerChainOption sets options.
type HandlerChainOption interface {
	apply(*HandlerChain)
}

// EmptyHandlerChainOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyHandlerChainOption struct{}

func (EmptyHandlerChainOption) apply(*HandlerChain) {}

// HandlerChainOptionFunc wraps a function that modifies HandlerChain into an
// implementation of the HandlerChainOption interface.
type HandlerChainOptionFunc func(*HandlerChain)

func (f HandlerChainOptionFunc) apply(do *HandlerChain) {
	f(do)
}

// sample code for option, default for nothing to change
func _HandlerChainOptionWithDefault() HandlerChainOption {
	return HandlerChainOptionFunc(func(*HandlerChain) {
		// nothing to change
	})
}
func (o *HandlerChain) ApplyOptions(options ...HandlerChainOption) *HandlerChain {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
