package cmd

import (
	"context"
	"fmt"
)

func runAllServices(_ context.Context, args *cmdArgs) error {
	if err := initLogger(args); err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}

	_, err := loadConfig(args)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	return fmt.Errorf("not implemented")
}
