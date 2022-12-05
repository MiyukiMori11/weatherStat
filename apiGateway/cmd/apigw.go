package cmd

import (
	"context"

	"github.com/MiyukiMori11/weatherstat/internal/config"
	"github.com/MiyukiMori11/weatherstat/internal/router"
	"github.com/MiyukiMori11/weatherstat/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const serviceName = "API Gateway"

var (
	RunApiGW = &cobra.Command{
		Use:  "run",
		RunE: ApiGatewayRunE,
	}
	cfgPath    string
	httpServer server.Server
)

func init() {
	RunApiGW.Flags().StringVarP(&cfgPath, "config", "", "./config/local.yaml", "Path to config file. Default: ./config/local.yaml")

	logger := zap.NewExample()
	logger = logger.Named(serviceName)

	cfg, err := config.Load(cfgPath)
	if err != nil {
		logger.Fatal("can't init config", zap.Error(err))
	}

	ginEngine := gin.Default()

	proxyRouter := router.New(&cfg.Gateway, logger, ginEngine)

	httpServer = server.New(
		&cfg.Server,
		ginEngine,
		logger,
		proxyRouter,
	)
}

func ApiGatewayRunE(command *cobra.Command, args []string) error {
	ctx, cancelFunc := context.WithCancel(command.Context())
	defer cancelFunc()

	if err := httpServer.Run(ctx); err != nil {
		return err
	}

	return nil
}
