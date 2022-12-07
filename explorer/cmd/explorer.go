package cmd

import (
	"fmt"
	"log"

	"github.com/MiyukiMori11/weatherstat/explorer/internal/client"
	"github.com/MiyukiMori11/weatherstat/explorer/internal/config"
	"github.com/MiyukiMori11/weatherstat/explorer/internal/parser"
	"github.com/MiyukiMori11/weatherstat/explorer/internal/server"
	"github.com/MiyukiMori11/weatherstat/explorer/internal/server/handler"
	"github.com/MiyukiMori11/weatherstat/explorer/internal/server/operations"
	"github.com/MiyukiMori11/weatherstat/explorer/internal/storage"
	"github.com/go-openapi/loads"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const serviceName = "explorer"

var (
	cfgPath string
	logger  *zap.Logger
)

func RunExplorerCommand() *cobra.Command {
	RunExplorer := &cobra.Command{
		Use:   "run",
		Short: "Runs explorer service",
		RunE:  ExplorerRunE,
	}

	RunExplorer.Flags().StringVar(&cfgPath, "config", "./config/local.yaml", "path to config file")

	return RunExplorer
}

func ExplorerRunE(command *cobra.Command, args []string) error {

	logger = zap.NewExample()
	logger = logger.Named(serviceName)

	cfg, err := config.Load(cfgPath)
	if err != nil {
		logger.Fatal("can't init config", zap.Error(err))
	}

	s, err := storage.New(command.Context(), cfg.Storage)
	if err != nil {
		return fmt.Errorf("error on storage init")
	}

	c := client.New(cfg.Client, logger)

	p := parser.New(cfg.Parser, logger, s, c)

	h := handler.New(logger, s, c)

	initServer(h)

	go p.Run(command.Context())

	return nil
}

func initServer(h handler.Handler) {
	swaggerSpec, err := loads.Embedded(server.SwaggerJSON, server.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}
	var s *server.Server

	api := operations.NewExplorerAPI(swaggerSpec)

	api.GetCitiesHandler = operations.GetCitiesHandlerFunc(h.GetCities)
	api.PostCitiesHandler = operations.PostCitiesHandlerFunc(h.PostCities)
	api.DeleteCitiesHandler = operations.DeleteCitiesHandlerFunc(h.DeleteCities)
	api.GetTempHandler = operations.GetTempHandlerFunc(h.GetTemp)

	server.Logger = logger

	s = server.NewServer(api)
	defer s.Shutdown()

	s.ConfigureAPI()
	if err := s.Serve(); err != nil {
		log.Fatalln(err)
	}
}
