package monitor

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"url-sentinel/internal/domain/entity"
	"url-sentinel/internal/domain/repository"
)

// Monitor periodically checks URLs and records results
type Monitor struct {
	urlRepo   repository.URLRepository
	checkRepo repository.CheckRepository
	client    *http.Client
	logger    *slog.Logger

	mu       sync.RWMutex
	watchers map[string]context.CancelFunc // urlID -> cancel function
}

// NewMonitor creates a new monitor instance
func NewMonitor(
	urlRepo repository.URLRepository,
	checkRepo repository.CheckRepository,
	logger *slog.Logger,
) *Monitor {
	return &Monitor{
		urlRepo:   urlRepo,
		checkRepo: checkRepo,
		client:    &http.Client{Timeout: 10 * time.Second},
		logger:    logger,
		watchers:  make(map[string]context.CancelFunc),
	}
}

// Start initializes monitoring for all URLs in the database
func (m *Monitor) Start(ctx context.Context) error {
	urls, err := m.urlRepo.List(ctx)
	if err != nil {
		m.logger.Error("failed to list urls for monitoring", slog.Any("error", err))
		return err
	}

	for _, url := range urls {
		m.AddURL(ctx, url)
	}

	m.logger.Info("monitor started", slog.Int("urls", len(urls)))
	return nil
}

// AddURL adds a new URL to monitoring
func (m *Monitor) AddURL(parentCtx context.Context, url *entity.URL) {
	m.mu.Lock()
	defer m.mu.Unlock()

	urlIDStr := url.ID.String()

	// Check if already monitoring
	if _, exists := m.watchers[urlIDStr]; exists {
		m.logger.Warn("url already being monitored", slog.String("url_id", urlIDStr))
		return
	}

	// Create cancellable context for this URL
	ctx, cancel := context.WithCancel(parentCtx)
	m.watchers[urlIDStr] = cancel

	// Start monitoring in a goroutine
	go m.watchURL(ctx, url)

	m.logger.Info("started monitoring url",
		slog.String("url_id", urlIDStr),
		slog.String("address", url.Address),
		slog.Duration("interval", url.CheckInterval),
	)
}

// RemoveURL stops monitoring a URL
func (m *Monitor) RemoveURL(urlID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cancel, exists := m.watchers[urlID]; exists {
		cancel()
		delete(m.watchers, urlID)
		m.logger.Info("stopped monitoring url", slog.String("url_id", urlID))
	}
}

// Stop stops all monitoring
func (m *Monitor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for urlID, cancel := range m.watchers {
		cancel()
		m.logger.Info("stopped monitoring url", slog.String("url_id", urlID))
	}

	m.watchers = make(map[string]context.CancelFunc)
	m.logger.Info("monitor stopped")
}

// watchURL performs periodic checks on a single URL
func (m *Monitor) watchURL(ctx context.Context, url *entity.URL) {
	ticker := time.NewTicker(url.CheckInterval)
	defer ticker.Stop()

	// Perform initial check immediately
	m.performCheck(ctx, url)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.performCheck(ctx, url)
		}
	}
}

// performCheck executes a single health check
func (m *Monitor) performCheck(ctx context.Context, url *entity.URL) {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.Address, nil)
	if err != nil {
		m.logger.Error("failed to create request",
			slog.String("url", url.Address),
			slog.Any("error", err),
		)
		return
	}

	resp, err := m.client.Do(req)
	duration := time.Since(start)

	status := false
	code := 0

	if err != nil {
		m.logger.Debug("check failed",
			slog.String("url", url.Address),
			slog.Any("error", err),
		)
	} else {
		defer resp.Body.Close()
		code = resp.StatusCode
		status = resp.StatusCode >= 200 && resp.StatusCode < 300
	}

	// Save check result
	check := entity.NewCheck(url.ID, status, code, duration)
	if err := m.checkRepo.Create(ctx, check); err != nil {
		m.logger.Error("failed to save check result",
			slog.String("url", url.Address),
			slog.Any("error", err),
		)
	} else {
		m.logger.Debug("check completed",
			slog.String("url", url.Address),
			slog.Int("code", code),
			slog.Bool("status", status),
			slog.Duration("duration", duration),
		)
	}
}
