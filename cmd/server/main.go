package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"url-sentinel/internal/config"
	"url-sentinel/internal/delivery/http/handler"
	mw "url-sentinel/internal/delivery/http/middleware"
	"url-sentinel/internal/monitor"
	"url-sentinel/internal/repository/postgres"
	"url-sentinel/internal/usecase"

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
	// Load configuration
	cfg := config.MustLoad()

	// Setup logger
	logger := setupLogger(cfg.Env)
	logger.Info("starting url-sentinel",
		slog.String("env", cfg.Env),
		slog.String("address", cfg.HTTPServer.Address),
	)

	// Initialize database
	db, err := postgres.New(cfg.Database.DSN())
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("failed to close database", slog.Any("error", err))
		}
	}()
	logger.Info("database connected successfully")

	// Initialize repositories
	urlRepo := postgres.NewURLRepository(db.DB)
	checkRepo := postgres.NewCheckRepository(db.DB)

	// Initialize and start monitor
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mon := monitor.NewMonitor(urlRepo, checkRepo, logger)
	if err := mon.Start(ctx); err != nil {
		logger.Error("failed to start monitor", slog.Any("error", err))
	}

	// Initialize use cases with monitor for dynamic URL management
	urlUseCase := usecase.NewURLUseCase(urlRepo, mon)
	checkUseCase := usecase.NewCheckUseCase(checkRepo)

	// Initialize handlers
	urlHandler := handler.NewURLHandler(urlUseCase, logger)
	checkHandler := handler.NewCheckHandler(checkUseCase, logger)

	// Setup router
	router := setupRouter(urlHandler, checkHandler, logger)

	// Setup HTTP server
	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("server listening", slog.String("address", cfg.HTTPServer.Address))
		serverErrors <- server.ListenAndServe()
	}()

	// Wait for interrupt signal or server error
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.Error("server error", slog.Any("error", err))
		os.Exit(1)

	case sig := <-shutdown:
		logger.Info("shutdown signal received", slog.String("signal", sig.String()))

		// Stop monitor
		mon.Stop()

		// Graceful shutdown with timeout
		ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPServer.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("failed to shutdown server gracefully", slog.Any("error", err))
			if err := server.Close(); err != nil {
				logger.Error("failed to close server", slog.Any("error", err))
			}
			os.Exit(1)
		}

		logger.Info("server stopped gracefully")
	}
}

func setupRouter(
	urlHandler *handler.URLHandler,
	checkHandler *handler.CheckHandler,
	logger *slog.Logger,
) *chi.Mux {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mw.Logger(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// URL routes
	router.Route("/urls", func(r chi.Router) {
		r.Post("/", urlHandler.Create)
		r.Get("/", urlHandler.List)
		r.Get("/{id}", urlHandler.Get)
		r.Delete("/{id}", urlHandler.Delete)
		r.Get("/{id}/history", checkHandler.GetHistory)
	})

	return router
}

func setupLogger(env string) *slog.Logger {
	var handler slog.Handler

	switch env {
	case envLocal, envDev:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	case envProd:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	return slog.New(handler)
}
