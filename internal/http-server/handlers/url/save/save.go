package save

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"log/slog"

	"url-sentinel/internal/model"
	"url-sentinel/internal/storage"

	"github.com/google/uuid"
)

// Request represents the payload for saving a new URL
type Request struct {
	Address       string `json:"address"`
	CheckInterval string `json:"check_interval"` // duration string, e.g. "30s"
}

// Response represents the saved URL data returned to the client
type Response struct {
	ID            uuid.UUID     `json:"id"`
	Address       string        `json:"address"`
	CheckInterval time.Duration `json:"check_interval"` // duration in nanoseconds
	CreatedAt     time.Time     `json:"created_at"`
}

// Handler registers the HTTP handler for saving a new URL
func Handler(repo *storage.URLRepository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("failed to decode request", slog.Any("error", err))
			http.Error(w, "invalid request payload", http.StatusBadRequest)
			return
		}

		// Parse interval from string
		interval, err := time.ParseDuration(req.CheckInterval)
		if err != nil {
			logger.Info("invalid check_interval format", slog.Any("error", err))
			http.Error(w, "invalid check_interval format", http.StatusBadRequest)
			return
		}

		u, err := model.NewURL(req.Address, interval)
		if err != nil {
			logger.Info("invalid URL parameters", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Save to repository
		if err := repo.SaveURL(u); err != nil {
			logger.Error("failed to save URL", slog.Any("error", err))
			if errors.Is(err, storage.ErrURLExists) {
				http.Error(w, "url already exists", http.StatusConflict)
			} else {
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Prepare response
		resp := Response{
			ID:            u.ID,
			Address:       u.Address,
			CheckInterval: u.CheckInterval,
			CreatedAt:     u.CreatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("failed to encode response", slog.Any("error", err))
		}
	}
}
