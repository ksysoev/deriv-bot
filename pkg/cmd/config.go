package cmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

type appConfig struct{}

// loadConfig loads the application configuration from the specified file path and environment variables.
// It uses the provided args structure to determine the configuration path.
// The function returns a pointer to the appConfig structure and an error if something goes wrong.
func loadConfig(args *cmdArgs) (*appConfig, error) {
	v := viper.NewWithOptions(viper.ExperimentalBindStruct())

	if args.ConfigPath != "" {
		v.SetConfigFile(args.ConfigPath)

		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var cfg appConfig

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	slog.Debug("Config loaded", slog.Any("config", cfg))

	return &cfg, nil
}
