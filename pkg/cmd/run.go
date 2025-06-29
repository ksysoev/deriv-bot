package cmd

import (
	"context"
	"fmt"

	"github.com/ksysoev/deriv-bot/pkg/core/executor"
	"github.com/ksysoev/deriv-bot/pkg/core/signal"
	"github.com/ksysoev/deriv-bot/pkg/prov/deriv"
)

func runAllServices(ctx context.Context, args *cmdArgs) error {
	if err := initLogger(args); err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}

	cfg, err := loadConfig(args)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	derivApi, err := deriv.New(cfg.Deriv)
	if err != nil {
		return fmt.Errorf("failed to create Deriv API client: %w", err)
	}

	defer derivApi.Close()

	marketSignals := signal.New(derivApi)

	exec := executor.New(marketSignals, derivApi)

	err = exec.ExecuteStrategy(ctx, args.Token, "R_100", 10, func(tick signal.Tick) bool {
		return true
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to ticks: %w", err)
	}

	return nil
}
