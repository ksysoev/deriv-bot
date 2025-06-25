package cmd

import (
	"github.com/spf13/cobra"
)

type BuildInfo struct {
	Version string
}

type cmdArgs struct {
	LogLevel   string `mapstructure:"log_level"`
	Version    string
	ConfigPath string `mapstructure:"config_path"`
	TextFormat bool   `mapstructure:"log_text"`
}

func InitCommand(build BuildInfo) cobra.Command {
	args := &cmdArgs{
		Version: build.Version,
	}

	cmd := cobra.Command{
		Use:   "bot",
		Short: "Service for running bots for Deriv API",
		Long:  "Service for running tranding bots and stream trading signal from markets data from Deriv API",
	}

	cmd.PersistentFlags().StringVar(&args.LogLevel, "log-level", "info", "log level (debug, info, warn, error)")
	cmd.PersistentFlags().BoolVar(&args.TextFormat, "log-text", true, "log in text format, otherwise JSON")
	cmd.PersistentFlags().StringVar(&args.ConfigPath, "config", "", "config path")

	cmd.AddCommand(initRunCommand(args))

	return cmd
}

func initRunCommand(args *cmdArgs) *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the bot",
		Long:  "Run the bot with the specified configuration.",
	}

	cmdRunAll := &cobra.Command{
		Use:   "all",
		Short: "Run all services",
		Long:  "Run all services",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runAllServices(cmd.Context(), args)
		},
	}

	runCmd.AddCommand(cmdRunAll)

	return runCmd
}
