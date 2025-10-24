package health

import (
	"io"
	"net/http"
	"proxymity/internal/backend"
	"proxymity/internal/config"
	"proxymity/internal/metrics"
	"time"
)

type HealthChecker struct {
	pool    *backend.Pool
	metrics *metrics.Metrics
	ticker  *time.Ticker
	timeout time.Duration
}

func NewHealthChecker(cfg config.HealthCheckConfig, pool *backend.Pool, m *metrics.Metrics) *HealthChecker {

	interval := time.Second * time.Duration(cfg.Interval)
	if cfg.Interval < 1 {
		interval = 5 * time.Second // Default interval to 5 seconds if not configured
	}

	timeout := time.Duration(cfg.TimeOut) * time.Second
	if cfg.TimeOut < 1 {
		timeout = 3 * time.Second // Default timeout to 3 seconds if not configured
	}

	return &HealthChecker{
		pool:    pool,
		metrics: m,
		ticker:  time.NewTicker(interval),
		timeout: timeout,
	}
}
func (h *HealthChecker) Start() {

	client := &http.Client{
		Timeout: h.timeout,
	}

	for range h.ticker.C {

		for _, b := range h.pool.GetBackends() {

			healthUrl := b.Host.JoinPath(b.Health)

			resp, err := client.Get(healthUrl.String())
			if err != nil {
				b.SetAlive(false)
				continue
			}

			// Read and discard the body to allow connection reuse, then close immediately.
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				b.SetAlive(true)
			} else {
				b.SetAlive(false)
			}
		}
	}
}

func (h *HealthChecker) Stop() {
	h.ticker.Stop()
}

func (h *HealthChecker) Backend(b *backend.Backend) bool {
	client := &http.Client{
		Timeout: h.timeout,
	}

	resp, err := client.Get(b.Host.Path)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
