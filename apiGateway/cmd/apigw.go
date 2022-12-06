package cmd

import (
	"context"

	"github.com/MiyukiMori11/weatherstat/apigateway/internal/config"
	"github.com/MiyukiMori11/weatherstat/apigateway/internal/router"
	"github.com/MiyukiMori11/weatherstat/apigateway/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const serviceName = "API Gateway"

var (
	cfgPath    string
	httpServer server.Server
)

func RunApiGWCommand() *cobra.Command {
	RunApiGW := &cobra.Command{
		Use:  "run",
		RunE: ApiGatewayRunE,
	}

	RunApiGW.Flags().StringVar(&cfgPath, "config", "./config/local.yaml", "path to config file")

	return RunApiGW
}

func ApiGatewayRunE(command *cobra.Command, args []string) error {

	initComponents()

	ctx, cancelFunc := context.WithCancel(command.Context())
	defer cancelFunc()

	if err := httpServer.Run(ctx); err != nil {
		return err
	}

	return nil
}

func initComponents() {
	logger := zap.NewExample()
	logger = logger.Named(serviceName)

	cfg, err := config.Load(cfgPath)
	if err != nil {
		logger.Fatal("can't init config", zap.Error(err))
	}

	ginEngine := gin.Default()

	proxyRouter := router.New(cfg.Gateway, logger, ginEngine)

	httpServer = server.New(
		cfg.Server,
		ginEngine,
		logger,
		proxyRouter,
	)
}
