package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/MiyukiMori11/weatherstat/internal/config"
	"github.com/MiyukiMori11/weatherstat/internal/proxy"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type server struct {
	config *config.Server
	engine *gin.Engine
	logger *zap.Logger

	proxy proxy.Proxy
}

type Server interface {
	Run(ctx context.Context) error
}

func New(
	config *config.Server,
	engine *gin.Engine,
	logger *zap.Logger,
	proxy proxy.Proxy) Server {
	return &server{
		config: config,
		engine: engine,
		logger: logger,
		proxy:  proxy,
	}
}

func (s *server) Run(ctx context.Context) error {
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	s.proxy.InitRoutes()

	server := &http.Server{
		Addr:    ":8080",
		Handler: s.engine,
	}

	go func() {
		defer cancelFunc()
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("listen error", zap.Error(err))
		}
	}()

	<-ctx.Done()

	s.logger.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancelFunc = context.WithTimeout(context.Background(), s.config.ShutdownTimeout())
	defer cancelFunc()
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	s.logger.Info("server exiting")

	return nil
}
