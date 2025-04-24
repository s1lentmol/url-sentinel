package main

import (
	"log/slog"
	"os"
	"url-sentinel/internal/config"
	"url-sentinel/internal/storage"

	_ "github.com/lib/pq"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting url-sentinel", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	db, err := storage.New(cfg.Database.DSN)
	if err != nil {
		log.Error("failed to open db", "error", err)
		os.Exit(1)
	}

	_ = db // delete

	log.Info("connected to database")

	// TODO: init router: chi

	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}
