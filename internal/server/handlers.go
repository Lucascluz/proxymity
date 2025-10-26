package server

import (
	"net/http"
	"proxymity/internal/backend"
	"proxymity/internal/config"
	"proxymity/internal/metrics"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// Health returns basic health status of the proxy
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "proxymity",
		"timestamp": time.Now().Unix(),
		"uptime":    time.Since(startTime).Seconds(),
	})
}

// Status returns detailed status including backend information
func Status(pool *backend.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		backends := pool.GetBackends()

		// Collect backend status
		backendStatus := make([]gin.H, 0, len(backends))
		healthyCount := 0

		for _, b := range backends {
			isHealthy := b.IsHealthy()
			if isHealthy {
				healthyCount++
			}

			backendStatus = append(backendStatus, gin.H{
				"name":    b.Name,
				"url":     b.Host.String(),
				"health":  b.Health,
				"enabled": b.IsEnabled(),
				"weight":  b.GetWeight(),
				"healthy": isHealthy,
			})
		}

		// Determine overall status
		var overallStatus string
		statusCode := http.StatusOK

		if healthyCount == 0 {
			overallStatus = "unhealthy"
			statusCode = http.StatusServiceUnavailable
		} else if healthyCount < len(backends) {
			overallStatus = "degraded"
		} else {
			overallStatus = "healthy"
		}

		// Get system stats
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		c.JSON(statusCode, gin.H{
			"status":    overallStatus,
			"service":   "proxymity",
			"timestamp": time.Now().Unix(),
			"uptime":    time.Since(startTime).Seconds(),
			"backends": gin.H{
				"total":   len(backends),
				"healthy": healthyCount,
				"details": backendStatus,
			},
			"system": gin.H{
				"goroutines":      runtime.NumGoroutine(),
				"memory_usage_mb": float64(m.Alloc) / 1024 / 1024,
				"cpu_count":       runtime.NumCPU(),
			},
		})
	}
}

// Config returns the current configuration settings
func Config(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":   "proxymity",
			"timestamp": time.Now().Unix(),
			"config":    config,
		})
	}
}

// Metrics returns the current metrics
func Metrics(metrics *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, metrics)
	}
}
