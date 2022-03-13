package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	http_ "github.com/kaydxh/golang/go/net/http"
)

type ProxyMode int32

const (
	Reverse_ProxyMode  ProxyMode = 0
	Forward_ProxyMode  ProxyMode = 1
	Redirect_ProxyMode ProxyMode = 2
)

type ProxyMatchedFunc func(c *gin.Context) bool
type ProxyTargetResolverFunc func(c *gin.Context) (host string, err error)

type proxyOptions struct {
	proxyMatchedFunc        ProxyMatchedFunc
	proxyTargetResolverFunc ProxyTargetResolverFunc
	proxyMode               ProxyMode
}

type Proxy struct {
	router gin.IRouter
	opts   proxyOptions
}

func defaultProxyOptions() proxyOptions {
	return proxyOptions{
		proxyMode: Reverse_ProxyMode,
	}
}

func NewProxy(router gin.IRouter, options ...ProxyOption) (*Proxy, error) {
	p := &Proxy{
		router: router,
	}
	p.opts = defaultProxyOptions()
	p.ApplyOptions(options...)

	p.setProxy()
	return p, nil
}

func (p *Proxy) ProxyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		if p.opts.proxyMatchedFunc != nil {
			// not apply proxy, process inplace
			if !p.opts.proxyMatchedFunc(c) {
				return
			}
		}

		if p.opts.proxyTargetResolverFunc == nil {
			// proxy target resolver func is nil, process inplace
			return
		}
		targetAddr, err := p.opts.proxyTargetResolverFunc(c)
		if err != nil {
			c.Render(http.StatusOK, render.JSON{Data: fmt.Errorf("resolve proxy target err: %v", err)})
			return
		}
		if targetAddr == "" {
			// proxy target is empty, process inplace
			return
		}

		targetUrl := http_.CloneURL(c.Request.URL)
		if targetUrl.Scheme == "" {
			targetUrl.Scheme = "http"
		}
		targetUrl.Host = targetAddr
		targetUrl.Path = "/"

		switch p.opts.proxyMode {
		case Redirect_ProxyMode:
			c.Redirect(http.StatusTemporaryRedirect, targetUrl.String())
			c.Abort()
			return

		case Reverse_ProxyMode:
			c.Request.Host = targetUrl.Host

		case Forward_ProxyMode:
		}

		rp := httputil.NewSingleHostReverseProxy(targetUrl)
		rp.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

func (p *Proxy) setProxy() {
	p.router.Use(p.ProxyHandler())
	/*
		for _, pattern := range p.opts.routerPatterns {
			p.router.Any(pattern, p.ProxyHandler())
		}
	*/
}
