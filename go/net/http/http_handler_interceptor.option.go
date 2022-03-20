package http

// A HandlerInterceptorsOption sets options.
type HandlerInterceptorsOption interface {
	apply(*HandlerInterceptors)
}

// EmptyHandlerInterceptorsOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyHandlerInterceptorsOption struct{}

func (EmptyHandlerInterceptorsOption) apply(*HandlerInterceptors) {}

// HandlerInterceptorsOptionFunc wraps a function that modifies HandlerInterceptors into an
// implementation of the HandlerInterceptorsOption interface.
type HandlerInterceptorsOptionFunc func(*HandlerInterceptors)

func (f HandlerInterceptorsOptionFunc) apply(do *HandlerInterceptors) {
	f(do)
}

// sample code for option, default for nothing to change
func _HandlerInterceptorsOptionWithDefault() HandlerInterceptorsOption {
	return HandlerInterceptorsOptionFunc(func(*HandlerInterceptors) {
		// nothing to change
	})
}
func (o *HandlerInterceptors) ApplyOptions(options ...HandlerInterceptorsOption) *HandlerInterceptors {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
