package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"net/http"
	"syscall"
	"url-sentinel/internal/config"
	"url-sentinel/internal/http-server/handlers/url/delete"
	"url-sentinel/internal/http-server/handlers/url/get"
	"url-sentinel/internal/http-server/handlers/url/history"
	"url-sentinel/internal/http-server/handlers/url/list"
	"url-sentinel/internal/http-server/handlers/url/save"
	mwLogger "url-sentinel/internal/http-server/middleware/logger"
	"url-sentinel/internal/monitor"
	"url-sentinel/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	// initialize repositories
	urlRepo := storage.NewURLRepository(db.GetDB())
	checkRepo := storage.NewCheckRepository(db.GetDB())

	// Start background monitor
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mon := monitor.NewMonitor(urlRepo, checkRepo, log)
	go mon.Start(ctx)

	// Listen for shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		log.Info("shutdown signal received")
		cancel()
	}()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.StripSlashes)

	router.Route("/urls", func(router chi.Router) {
		router.Post("/", save.Handler(urlRepo, log))
		router.Get("/", list.Handler(urlRepo, log))
		router.Get("/{id}", get.Handler(urlRepo, log))
		router.Delete("/{id}", delete.Handler(urlRepo, log))
		router.Get("/{id}/history", history.Handler(checkRepo, log))
	})
	log.Info("listening on address", slog.String("addr", cfg.HTTPServer.Address))
	if err := http.ListenAndServe(cfg.HTTPServer.Address, router); err != nil {
		log.Error("server error", slog.Any("error", err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	// can add different environmen
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}
