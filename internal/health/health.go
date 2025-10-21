package health

import (
	"io"
	"net/http"
	"proxymity/internal/backend"
	"proxymity/internal/config"
	"time"
)

type HealthCheck struct {
	pool     *backend.Pool
	interval *time.Ticker
	timeout  time.Duration
}

func NewHealthCheck(cfg config.HealthCheckConfig, pool *backend.Pool) *HealthCheck {

	interval := time.Second * time.Duration(cfg.Interval)
	if cfg.Interval < 1 {
		interval = 5 * time.Second // Default interval to 5 seconds if not configured
	}

	timeout := time.Duration(cfg.TimeOut) * time.Second
	if cfg.TimeOut < 1 {
		timeout = 3 * time.Second // Default timeout to 3 seconds if not configured
	}

	return &HealthCheck{
		pool:     pool,
		interval: time.NewTicker(interval),
		timeout:  timeout,
	}
}
func (h *HealthCheck) Start() {

	client := &http.Client{
		Timeout: h.timeout,
	}

	for range h.interval.C {

		for _, b := range h.pool.GetBackends() {

			resp, err := client.Get(b.URL.String())
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

func (h *HealthCheck) Stop() {
	h.interval.Stop()
}

func (h *HealthCheck) Backend(b *backend.Backend) bool {
	client := &http.Client{
		Timeout: h.timeout,
	}

	resp, err := client.Get(b.URL.Path)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
