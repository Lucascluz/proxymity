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

func NewHealthChecker(cfg config.HealthCheckConfig, pool *backend.Pool, metrics *metrics.Metrics) *HealthChecker {

	interval := time.Second * time.Duration(cfg.Interval)
	ticker := time.NewTicker(interval)
	timeout := time.Duration(cfg.TimeOut) * time.Second

	return &HealthChecker{
		pool:    pool,
		metrics: metrics,
		ticker:  ticker,
		timeout: timeout,
	}
}
func (hc *HealthChecker) Start() {

	client := &http.Client{
		Timeout: hc.timeout,
	}

	for range hc.ticker.C {

		for _, b := range hc.pool.GetBackends() {

			// Skip health check if in backoff period
			if time.Since(b.GetChecked()) < b.GetBackoff() {
				continue
			}

			// Health check
			b.SetChecked()
			healthUrl := b.Host.JoinPath(b.Health)
			resp, err := client.Get(healthUrl.String())
			if err != nil {
				hc.metrics.Error.IncClientErrs()
				b.SetHealthy(false) // Mark unhealthy on network failure
				b.ExpBackof()       // Increase backoff
				continue
			}

			// Read and discard the body to allow connection reuse, then close immediately.
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				b.SetHealthy(true)
				b.ResetBackof()
			} else {
				b.SetHealthy(false)
				b.ExpBackof()
			}
		}
	}
}

func (hc *HealthChecker) Stop() {
	hc.ticker.Stop()
}
