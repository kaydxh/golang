package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	http_ "github.com/kaydxh/golang/go/net/http"
)

type ProxyMode int32

const (
	Reverse_ProxyMode  ProxyMode = 0
	Forward_ProxyMode  ProxyMode = 1
	Redirect_ProxyMode ProxyMode = 2
)

type MatchRouterFunc func(*http.Request) string

type Proxy struct {
	router gin.IRouter
	opts   struct {
		routerPatterns []string
		matchRouter    MatchRouterFunc
		targetUrl      string
		proxyMode      ProxyMode
	}
}

func NewProxy(router gin.IRouter, options ...ProxyOption) (*Proxy, error) {
	p := &Proxy{
		router: router,
	}
	p.ApplyOptions(options...)
	if p.opts.targetUrl == "" && p.opts.matchRouter == nil {
		return nil, fmt.Errorf("target url and match router both nil")
	}

	return p, nil
}

func (p *Proxy) ProxyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request

		serviceTargetUrl := http_.CloneURL(c.Request.URL)
		if serviceTargetUrl.Scheme == "" {
			serviceTargetUrl.Scheme = "http"
		}
		serviceTargetUrl.Path = "/"
		sTargetUrl := p.opts.targetUrl
		if p.opts.matchRouter != nil {
			sTargetUrl = p.opts.matchRouter(req)
		}

		if sTargetUrl == "" {
			return
		}
		targetUrl, err := url.Parse(sTargetUrl)
		if err != nil {
			return
		}
		if targetUrl.Host != "" {
			serviceTargetUrl.Host = targetUrl.Host
		}

		switch p.opts.proxyMode {
		case Redirect_ProxyMode:
			c.Redirect(http.StatusTemporaryRedirect, serviceTargetUrl.String())
			c.Abort()
			return

		case Reverse_ProxyMode:
			c.Request.Host = serviceTargetUrl.Host

		case Forward_ProxyMode:
		}

		rp := httputil.NewSingleHostReverseProxy(serviceTargetUrl)
		rp.ServeHTTP(c.Writer, c.Request)
	}
}

func (p *Proxy) SetProxy() {
	for _, pattern := range p.opts.routerPatterns {
		p.router.Any(pattern, p.ProxyHandler())
	}
}
