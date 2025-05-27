package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"url-sentinel/internal/delivery/http/dto"
	"url-sentinel/internal/usecase"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// CheckHandler handles HTTP requests for check operations
type CheckHandler struct {
	checkUseCase *usecase.CheckUseCase
	logger       *slog.Logger
}

// NewCheckHandler creates a new check handler
func NewCheckHandler(checkUseCase *usecase.CheckUseCase, logger *slog.Logger) *CheckHandler {
	return &CheckHandler{
		checkUseCase: checkUseCase,
		logger:       logger,
	}
}

// GetHistory handles GET /urls/{id}/history
func (h *CheckHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Info("invalid url id", slog.String("id", idParam))
		h.respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	checks, err := h.checkUseCase.GetCheckHistory(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get check history", slog.Any("error", err))
		h.respondError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := make([]dto.CheckResponse, 0, len(checks))
	for _, check := range checks {
		resp = append(resp, dto.CheckResponse{
			ID:        check.ID,
			URLID:     check.URLID,
			Status:    check.Status,
			Code:      check.Code,
			Duration:  check.Duration.String(),
			CheckedAt: check.CheckedAt,
		})
	}

	h.respondJSON(w, resp, http.StatusOK)
}

func (h *CheckHandler) respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode response", slog.Any("error", err))
	}
}

func (h *CheckHandler) respondError(w http.ResponseWriter, message string, status int) {
	h.respondJSON(w, dto.ErrorResponse{Error: message}, status)
}
