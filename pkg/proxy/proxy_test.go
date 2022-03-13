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
