package get

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	"url-sentinel/internal/storage"
)

// Handler returns a specific URL by ID
func Handler(repo *storage.URLRepository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			logger.Info("invalid URL id", slog.String("id", idParam), slog.Any("error", err))
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		u, err := repo.GetURLByID(id)
		if err != nil {
			logger.Error("failed to get URL", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if u == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(u); err != nil {
			logger.Error("failed to encode response", slog.Any("error", err))
		}
	}
}
