package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/MiyukiMori11/weatherstat/apigateway/cmd"
	"github.com/spf13/cobra"
)

func main() {
	ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelFunc()

	go func() {
		<-ctx.Done()
		cancelFunc()
	}()

	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(cmd.RunApiGWCommand())
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
