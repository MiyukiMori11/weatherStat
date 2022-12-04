package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/MiyukiMori11/weatherstat/internal/config"
	"github.com/MiyukiMori11/weatherstat/internal/domain"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Proxy interface {
	InitRoutes()
}

type proxy struct {
	logger *zap.Logger
	config *config.Gateway

	engine *gin.Engine
}

func New(config *config.Gateway, logger *zap.Logger, engine *gin.Engine) Proxy {
	return &proxy{
		config: config,
		logger: logger,
		engine: engine,
	}
}

func (p *proxy) InitRoutes() {
	for _, service := range p.config.Services {
		rGroup := p.engine.Group(service.Root)
		for _, route := range service.Routes {
			handler := p.proxyHandler(
				service.Scheme+"://"+service.Host+":"+service.Port,
				route.Path)
			switch {
			case route.Method == domain.GET:
				rGroup.GET(route.Path, handler)
			case route.Method == domain.POST:
				rGroup.POST(route.Path, handler)
			case route.Method == domain.DELETE:
				rGroup.DELETE(route.Path, handler)
			case route.Method == domain.PUT:
				rGroup.PUT(route.Path, handler)
			}
		}
	}

}

func (p *proxy) proxyHandler(serviceUrl, path string) gin.HandlerFunc {
	remote, err := url.Parse(serviceUrl)
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		proxy := httputil.NewSingleHostReverseProxy(remote)

		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = path
			req.Body = c.Request.Body
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
