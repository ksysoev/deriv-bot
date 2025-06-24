package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ksysoev/deriv-bot/pkg/cmd"
)

var version = "dev"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	command := cmd.InitCommand(cmd.BuildInfo{
		Version: version,
	})

	err := command.ExecuteContext(ctx)
	if err != nil {
		fmt.Println(err)
		cancel()

		os.Exit(1)
	}

	cancel()
}
