/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
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
			// filter interceptor is nil
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
