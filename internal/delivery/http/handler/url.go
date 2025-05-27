package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"url-sentinel/internal/delivery/http/dto"
	"url-sentinel/internal/domain/repository"
	"url-sentinel/internal/usecase"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// URLHandler handles HTTP requests for URL operations
type URLHandler struct {
	urlUseCase *usecase.URLUseCase
	logger     *slog.Logger
}

// NewURLHandler creates a new URL handler
func NewURLHandler(urlUseCase *usecase.URLUseCase, logger *slog.Logger) *URLHandler {
	return &URLHandler{
		urlUseCase: urlUseCase,
		logger:     logger,
	}
}

// Create handles POST /urls
func (h *URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", slog.Any("error", err))
		h.respondError(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	// Parse interval
	interval, err := time.ParseDuration(req.CheckInterval)
	if err != nil {
		h.logger.Info("invalid check_interval format", slog.Any("error", err))
		h.respondError(w, "invalid check_interval format", http.StatusBadRequest)
		return
	}

	// Create URL
	url, err := h.urlUseCase.CreateURL(r.Context(), req.Address, interval)
	if err != nil {
		if errors.Is(err, repository.ErrURLAddressExists) {
			h.logger.Info("url already exists", slog.String("address", req.Address))
			h.respondError(w, "url already exists", http.StatusConflict)
			return
		}
		h.logger.Error("failed to create url", slog.Any("error", err))
		h.respondError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Prepare response
	resp := dto.URLResponse{
		ID:            url.ID,
		Address:       url.Address,
		CheckInterval: url.CheckInterval.String(),
		CreatedAt:     url.CreatedAt,
	}

	h.respondJSON(w, resp, http.StatusCreated)
}

// Get handles GET /urls/{id}
func (h *URLHandler) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Info("invalid url id", slog.String("id", idParam))
		h.respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	url, err := h.urlUseCase.GetURLByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrURLNotFound) {
			h.respondError(w, "url not found", http.StatusNotFound)
			return
		}
		h.logger.Error("failed to get url", slog.Any("error", err))
		h.respondError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := dto.URLResponse{
		ID:            url.ID,
		Address:       url.Address,
		CheckInterval: url.CheckInterval.String(),
		CreatedAt:     url.CreatedAt,
	}

	h.respondJSON(w, resp, http.StatusOK)
}

// List handles GET /urls
func (h *URLHandler) List(w http.ResponseWriter, r *http.Request) {
	urls, err := h.urlUseCase.ListURLs(r.Context())
	if err != nil {
		h.logger.Error("failed to list urls", slog.Any("error", err))
		h.respondError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := make([]dto.URLResponse, 0, len(urls))
	for _, url := range urls {
		resp = append(resp, dto.URLResponse{
			ID:            url.ID,
			Address:       url.Address,
			CheckInterval: url.CheckInterval.String(),
			CreatedAt:     url.CreatedAt,
		})
	}

	h.respondJSON(w, resp, http.StatusOK)
}

// Delete handles DELETE /urls/{id}
func (h *URLHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Info("invalid url id", slog.String("id", idParam))
		h.respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.urlUseCase.DeleteURL(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrURLNotFound) {
			h.respondError(w, "url not found", http.StatusNotFound)
			return
		}
		h.logger.Error("failed to delete url", slog.Any("error", err))
		h.respondError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *URLHandler) respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode response", slog.Any("error", err))
	}
}

func (h *URLHandler) respondError(w http.ResponseWriter, message string, status int) {
	h.respondJSON(w, dto.ErrorResponse{Error: message}, status)
}
