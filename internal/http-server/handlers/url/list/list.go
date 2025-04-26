package list

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"url-sentinel/internal/storage"
)

// Handler returns all monitored URLs
func Handler(repo *storage.URLRepository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urls, err := repo.ListOfURLs()
		if err != nil {
			logger.Error("failed to list URLs", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(urls); err != nil {
			logger.Error("failed to encode response", slog.Any("error", err))
		}
	}
}
