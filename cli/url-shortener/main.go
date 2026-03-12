package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/RishatShay/url-shortener/internal/config"
	"github.com/RishatShay/url-shortener/internal/storage/sqlite"
	"github.com/RishatShay/url-shortener/internal/utils/logger"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Init config
	cfg := config.MustLoad()

	// Init logger
	log := MustSetupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Init storage
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog.String("error", err.Error()))
		os.Exit(1)
	}
	_ = storage

	// Init router
	router := chi.NewRouter()
	router.Use(logger.New(log))

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
