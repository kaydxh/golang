package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type MatchRouterFunc func(*http.Request) string

type ReverseProxy struct {
	router     gin.IRouter
	targetUrls []string
	opts       struct {
		routerPatterns []string
		matchRouter    MatchRouterFunc
	}
}

func NewReverseProxy(router gin.IRouter, targetUrls []string, options ...ReverseProxyOption) *ReverseProxy {
	p := &ReverseProxy{
		router:     router,
		targetUrls: targetUrls,
	}
	p.ApplyOptions(options...)

	return p
}

func (p *ReverseProxy) ProxyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("--ProxyHandler")
		req := c.Request

		var targetUrl string
		if p.opts.matchRouter != nil {
			targetUrl = p.opts.matchRouter(req)
		}

		if targetUrl == "" {
			return
		}

		serviceTargetUrl, err := url.Parse(targetUrl)
		if err != nil {
			return
		}
		rp := httputil.NewSingleHostReverseProxy(serviceTargetUrl)
		rp.ServeHTTP(c.Writer, c.Request)
	}

}

func (p *ReverseProxy) SetProxy() {
	for _, pattern := range p.opts.routerPatterns {
		p.router.Any(pattern, p.ProxyHandler())
	}
}
