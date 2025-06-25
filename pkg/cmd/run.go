package cmd

import (
	"context"
	"fmt"
	"log/slog"

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

	sub, err := derivApi.SubscribeToTicks(ctx, "R_100")
	if err != nil {
		return fmt.Errorf("failed to subscribe to ticks: %w", err)
	}

	for tick := range sub {
		slog.Info("Received tick", slog.Any("tick", tick))
	}

	return nil
}
