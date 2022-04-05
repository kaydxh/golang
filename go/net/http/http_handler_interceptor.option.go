package http

// A HandlerInterceptorOption sets options.
type HandlerInterceptorOption interface {
	apply(*HandlerInterceptor)
}

// EmptyHandlerInterceptorOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyHandlerInterceptorOption struct{}

func (EmptyHandlerInterceptorOption) apply(*HandlerInterceptor) {}

// HandlerInterceptorOptionFunc wraps a function that modifies HandlerInterceptor into an
// implementation of the HandlerInterceptorOption interface.
type HandlerInterceptorOptionFunc func(*HandlerInterceptor)

func (f HandlerInterceptorOptionFunc) apply(do *HandlerInterceptor) {
	f(do)
}

// sample code for option, default for nothing to change
func _HandlerInterceptorOptionWithDefault() HandlerInterceptorOption {
	return HandlerInterceptorOptionFunc(func(*HandlerInterceptor) {
		// nothing to change
	})
}
func (o *HandlerInterceptor) ApplyOptions(options ...HandlerInterceptorOption) *HandlerInterceptor {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
