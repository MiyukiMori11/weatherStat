// Package cmd migrate represents migration command
package cmd

import (
	"database/sql"
	"errors"

	"github.com/MiyukiMori11/weatherstat/explorer/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func MigrateCommand() *cobra.Command {
	RunMigration := &cobra.Command{
		Use:   "migrate",
		Short: "Runs migrations for explorer service",
		RunE:  MigrateE,
	}

	return RunMigration
}

func MigrateE(command *cobra.Command, args []string) error {

	c, err := config.Load(cfgPathExplorer)
	if err != nil {
		logger.Fatal("can't init config", zap.Error(err))
	}

	cfg := c.Storage

	db, err := sql.Open("postgres", cfg.URL()+"?sslmode=disable")
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	return nil
}
