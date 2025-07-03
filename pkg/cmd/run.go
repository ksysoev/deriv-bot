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

	initPrice := float64(0)
	strategy := executor.Strategy{
		Token:    args.Token,
		Symbol:   "R_100",
		Amount:   10,
		Type:     executor.StrategyTypeBuy,
		Leverage: 10,
		CheckToOpen: func(tick signal.Tick) bool {
			initPrice = tick.Quote
			return true
		},
		CheckToClose: func(tick signal.Tick) bool {
			return tick.Quote > initPrice*1.01
		},
	}

	err = exec.ExecuteStrategy(ctx, strategy)
	if err != nil {
		return fmt.Errorf("failed to subscribe to ticks: %w", err)
	}

	return nil
}
