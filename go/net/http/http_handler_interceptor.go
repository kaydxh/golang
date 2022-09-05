package http

import (
	"net/http"

	runtime_ "github.com/kaydxh/golang/go/runtime"
)

type HandlerInterceptor struct {
	//http middleware
	Interceptor func(h http.Handler) http.Handler
}

func NewHandlerInterceptor(opts ...HandlerInterceptorOption) *HandlerInterceptor {
	handers := &HandlerInterceptor{}
	handers.ApplyOptions(opts...)
	return handers
}

type HandlerChain struct {
	//invoke before http handler
	PreHandlers []func(w http.ResponseWriter, r *http.Request) error
	Handlers    []HandlerInterceptor
	//invoke after http handler
	PostHandlers []func(w http.ResponseWriter, r *http.Request)
}

func NewHandlerChain(opts ...HandlerChainOption) *HandlerChain {
	c := &HandlerChain{}
	c.ApplyOptions(opts...)

	return c
}

func (c *HandlerChain) WrapH(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer runtime_.Recover()

		for _, preH := range c.PreHandlers {
			err := preH(w, r)
			if err != nil {
				// assume PreHandler already process response by self
				return
			}
		}

		// reverse iterate handlers, so called handler as registed order
		for i := len(c.Handlers) - 1; i >= 0; i-- {
			if c.Handlers[i].Interceptor != nil {
				next = c.Handlers[i].Interceptor(next)
			}
		}

		next.ServeHTTP(w, r)

		for _, postH := range c.PostHandlers {
			postH(w, r)
		}
	})
}
