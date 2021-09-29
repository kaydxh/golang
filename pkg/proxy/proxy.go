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
	router gin.IRouter
	opts   struct {
		routerPatterns []string
		matchRouter    MatchRouterFunc
		targetUrl      string
	}
}

func NewReverseProxy(router gin.IRouter, options ...ReverseProxyOption) (*ReverseProxy, error) {
	p := &ReverseProxy{
		router: router,
	}
	p.ApplyOptions(options...)
	if p.opts.targetUrl == "" && p.opts.matchRouter == nil {
		return nil, fmt.Errorf("target url and match router both nil")
	}

	return p, nil
}

func (p *ReverseProxy) ProxyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request

		targetUrl := p.opts.targetUrl
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
