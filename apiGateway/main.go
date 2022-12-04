package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/MiyukiMori11/weatherstat/cmd"
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
	rootCmd.AddCommand(cmd.RunApiGW)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
