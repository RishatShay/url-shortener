package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/RishatShay/url-shortener/internal/config"
)

func main() {
	// Init config
	cfg := config.MustLoad()

	// Init logger
	logger := MustSetupLogger(cfg.Env)
	logger.Info("starting url-shortener", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

	// Init storage

	// Init router

	// Run server
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func MustSetupLogger(env string) *slog.Logger {
	var loger *slog.Logger

	switch env {
	case envLocal:
		loger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		loger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		loger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log.Fatalf("unable to create loger: incorrect environment (env) value: %s", env)
	}

	return loger
}
