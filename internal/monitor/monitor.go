package monitor

import (
	"context"
	"net/http"
	"time"

	"log/slog"

	"url-sentinel/internal/model"
	"url-sentinel/internal/storage"
)

// Monitor periodically checks URLs and records results
type Monitor struct {
	urlRepo   *storage.URLRepository
	checkRepo *storage.CheckRepository
	client    *http.Client
	logger    *slog.Logger
}

// NewMonitor constructs a Monitor
func NewMonitor(
	urlRepo *storage.URLRepository,
	checkRepo *storage.CheckRepository,
	logger *slog.Logger,
) *Monitor {
	return &Monitor{
		urlRepo:   urlRepo,
		checkRepo: checkRepo,
		client:    &http.Client{Timeout: 10 * time.Second},
		logger:    logger,
	}
}

// Start begins monitoring all URLs in background
func (m *Monitor) Start(ctx context.Context) {
	urls, err := m.urlRepo.ListOfURLs()
	if err != nil {
		m.logger.Error("monitor: failed to list URLs", slog.Any("error", err))
		return
	}
	for _, u := range urls {
		// spawn a goroutine per URL
		go m.watchURL(ctx, u)
	}
}

// watchURL periodically checks a single URL
func (m *Monitor) watchURL(ctx context.Context, u model.URL) {
	ticker := time.NewTicker(u.CheckInterval)
	defer ticker.Stop()

	m.logger.Info("monitoring URL", slog.String("url", u.Address), slog.Duration("interval", u.CheckInterval))
	// initial check
	m.doCheck(u)

	for {
		select {
		case <-ctx.Done():
			m.logger.Info("stopping monitor for URL", slog.String("url", u.Address))
			return
		case <-ticker.C:
			m.doCheck(u)
		}
	}
}

// doCheck performs a single HTTP GET and records the result
func (m *Monitor) doCheck(u model.URL) {
	t0 := time.Now()
	resp, err := m.client.Get(u.Address)
	duration := time.Since(t0)
	status := false
	code := 0
	if err != nil {
		m.logger.Error("monitor: request failed", slog.String("url", u.Address), slog.Any("error", err))
	} else {
		defer resp.Body.Close()
		code = resp.StatusCode
		status = resp.StatusCode < 300
	}

	check := model.NewCheck(u.ID, status, code, duration)
	if err := m.checkRepo.SaveCheck(check); err != nil {
		m.logger.Error("monitor: failed to save check", slog.Any("error", err))
	}
}
