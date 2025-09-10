package main

import (
	"log/slog"
	"os"

	"github.com/svetsed/todo_cli_app/cmd"
	"github.com/svetsed/todo_cli_app/internal/config"
	"github.com/svetsed/todo_cli_app/internal/logger"
)

func main() {
	logger.Init(slog.LevelError, os.Stdout)

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("failed to load configuration", err)
		os.Exit(1)
	}

	rootCmd := cmd.RootCmd(cfg)
	if err := rootCmd.Execute(); err != nil {
		logger.Error("failed to execute root command", err)
		os.Exit(1)
	}
}
