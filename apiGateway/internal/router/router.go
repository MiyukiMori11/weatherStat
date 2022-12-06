package router

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/MiyukiMori11/weatherstat/apigateway/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Router interface {
	InitRoutes() error
}

type router struct {
	logger *zap.Logger
	config *config.Gateway

	engine *gin.Engine
}

func New(config *config.Gateway, logger *zap.Logger, engine *gin.Engine) Router {
	return &router{
		config: config,
		logger: logger,
		engine: engine,
	}
}

func (r *router) InitRoutes() error {
	for _, service := range r.config.Services {
		rGroup := r.engine.Group(service.Root)
		for _, route := range service.Routes {
			handler, err := r.proxyHandler(
				service.Scheme+"://"+service.Host+":"+service.Port,
				route.Path)
			if err != nil {
				return fmt.Errorf("can't create handler for %s: %w", service.Name, err)
			}

			switch {
			case route.Method == GET:
				rGroup.GET(route.Path, handler)
			case route.Method == POST:
				rGroup.POST(route.Path, handler)
			case route.Method == DELETE:
				rGroup.DELETE(route.Path, handler)
			case route.Method == PUT:
				rGroup.PUT(route.Path, handler)
			}
		}
	}

	return nil

}

func (r *router) proxyHandler(serviceUrl, path string) (gin.HandlerFunc, error) {
	remote, err := url.Parse(serviceUrl)
	if err != nil {
		return nil, fmt.Errorf("can't parse url: %w", err)
	}

	return func(c *gin.Context) {
		router := httputil.NewSingleHostReverseProxy(remote)

		router.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = path
			req.Body = c.Request.Body
		}

		router.ServeHTTP(c.Writer, c.Request)
	}, nil
}
