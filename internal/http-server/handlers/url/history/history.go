package history

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	"url-sentinel/internal/storage"
)

// Handler returns check history for a URL
func Handler(checkRepo *storage.CheckRepository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			logger.Info("invalid URL id for history", slog.String("id", idParam), slog.Any("error", err))
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		checks, err := checkRepo.ListOfChecksByURL(id)
		if err != nil {
			logger.Error("failed to get check history", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(checks); err != nil {
			logger.Error("failed to encode response", slog.Any("error", err))
		}
	}
}
