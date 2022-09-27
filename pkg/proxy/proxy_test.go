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
package proxy_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	proxy_ "github.com/kaydxh/golang/pkg/proxy"
)

func TestNewProxy(t *testing.T) {
	type args struct {
		router  gin.IRouter
		options []proxy_.ProxyOption
	}
	r := gin.Default()
	tests := []struct {
		name    string
		args    args
		want    *proxy_.Proxy
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "proxy",
			args: args{
				router: r,
				options: []proxy_.ProxyOption{
					proxy_.WithProxyTargetResolverFunc(
						func(c *gin.Context) (string, error) {
							return "127.0.0.1:8081", nil
						},
					),
					//	proxy_.WithProxyMode(proxy_.Redirect_ProxyMode),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proxy, err := proxy_.NewProxy(tt.args.router, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProxy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Use attaches a global middleware to the router. ie. the middleware attached though Use() will be
			// included in the handlers chain for every single request. Even 404, 405, static files...
			// For example, this is the right place for a logger or error management middleware.
			r.Use(proxy.ProxyHandler())
			r.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})
			r.Run()
		})
	}
}
