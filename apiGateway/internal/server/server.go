package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/MiyukiMori11/weatherstat/apigateway/internal/config"
	"github.com/MiyukiMori11/weatherstat/apigateway/internal/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type server struct {
	config *config.Server
	engine *gin.Engine
	logger *zap.Logger

	router router.Router
}

type Server interface {
	Run(ctx context.Context) error
}

func New(
	config *config.Server,
	engine *gin.Engine,
	logger *zap.Logger,
	router router.Router) Server {
	return &server{
		config: config,
		engine: engine,
		logger: logger,
		router: router,
	}
}

func (s *server) Run(ctx context.Context) error {
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	if err := s.router.InitRoutes(); err != nil {
		return fmt.Errorf("can't init routes: %w", err)
	}

	server := &http.Server{
		Addr:    s.config.Addr(),
		Handler: s.engine,
	}

	go func() {
		defer cancelFunc()
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("listen error", zap.Error(err))
			os.Exit(1)
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
