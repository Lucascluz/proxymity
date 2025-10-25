package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	loadbalancer "proxymity/internal/balancer"
	"proxymity/internal/metrics"

	"github.com/gin-gonic/gin"
)

type Proxy struct {
	lb loadbalancer.LoadBalancer
	m  *metrics.Metrics
}

func NewProxy(lb loadbalancer.LoadBalancer, m *metrics.Metrics) *Proxy {
	return &Proxy{lb: lb}
}

func (p *Proxy) Proxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			lastErr  error
			tried    int
			maxTries = p.lb.CountAvailableBackends()
		)
		for tried < maxTries {
			backend, err := p.lb.NextBackend()
			if err != nil {
				lastErr = err
				break
			}

			proxy := httputil.NewSingleHostReverseProxy(backend.Host)
			proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
				log.Printf("Error proxying to %s: %v", backend.Name, err)
				backend.SetAlive(false)
				lastErr = err
			}

			proxy.ServeHTTP(c.Writer, c.Request)
			backend.AddConnection()

			// If no error was set by ErrorHandler, request succeeded
			if lastErr == nil {
				return
			}

			tried++
		}
		// If we reach here, all attempts failed
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "all backends failed",
			"details": lastErr.Error(),
		})
	}
}
