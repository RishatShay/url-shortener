package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/RishatShay/url-shortener/internal/config"
	"github.com/RishatShay/url-shortener/internal/http-server/handlers/url/save"
	"github.com/RishatShay/url-shortener/internal/http-server/middleware/logger"
	"github.com/RishatShay/url-shortener/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	// Run server

	log.Info("starting server", "address", cfg.HTTPServer.Address)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.Idle_timeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server is down")
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
