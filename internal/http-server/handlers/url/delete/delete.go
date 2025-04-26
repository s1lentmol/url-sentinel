package delete

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	"url-sentinel/internal/storage"
)

// Handler removes a URL by ID
func Handler(repo *storage.URLRepository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			logger.Info("invalid URL id for delete", slog.String("id", idParam), slog.Any("error", err))
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := repo.Delete(id); err != nil {
			logger.Error("failed to delete URL", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}