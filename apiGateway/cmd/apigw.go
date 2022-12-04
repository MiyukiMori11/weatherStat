package cmd

import (
	"context"

	"github.com/MiyukiMori11/weatherstat/internal/config"
	"github.com/MiyukiMori11/weatherstat/internal/proxy"
	"github.com/MiyukiMori11/weatherstat/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var RunApiGW = &cobra.Command{
	Use:  "gateway",
	RunE: ApiGatewayRunE,
}

func ApiGatewayRunE(command *cobra.Command, args []string) error {
	ctx, cancelFunc := context.WithCancel(command.Context())
	defer cancelFunc()

	logger := zap.NewExample()
	logger = logger.Named(command.Name())

	cfg, err := config.Load("./config")
	if err != nil {
		logger.Fatal("can't init config", zap.Error(err))
	}

	engine := gin.Default()

	p := proxy.New(&cfg.Gateway, logger, engine)

	s := server.New(
		&cfg.Server,
		engine,
		logger,
		p,
	)

	s.Run(ctx)

	return nil
}
