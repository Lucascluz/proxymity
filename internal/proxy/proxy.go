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
	return &Proxy{lb: lb,
		m: m}
}

// responseWriter wraps gin.ResponseWriter to track bytes written
type responseWriter struct {
	gin.ResponseWriter
	bytesWritten int64
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(data)
	rw.bytesWritten += int64(n)
	return n, err
}

func (p *Proxy) Proxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			lastErr  error
			tried    int
			maxTries = p.lb.CountAvailableBackends()
		)

		log.Printf("\n %d available backends \n", maxTries)

		for tried < maxTries {
			backend, err := p.lb.NextBackend()
			if err != nil {
				lastErr = err
				break
			}

			proxy := httputil.NewSingleHostReverseProxy(backend.Host)

			// Customize the error handler to capture errors
			proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
				log.Printf("Error proxying to %s: %v", backend.Name, err)
				lastErr = err
				p.m.Error.IncServerErrs()
			}

			// Wrap the response writer to track bytes out
			rw := &responseWriter{ResponseWriter: c.Writer}

			proxy.ServeHTTP(rw, c.Request)
			backend.AddConnection()

			// Update metrics
			p.m.Traffic.IncRequests()
			p.m.Traffic.AddBytesIn(c.Request.ContentLength)
			p.m.Traffic.AddBytesOut(rw.bytesWritten)

			// If no error was set by ErrorHandler, request succeeded
			if lastErr == nil {
				return
			}

			tried++
		}
		// If all attempts failed
		details := "no backends available"
		if lastErr != nil {
			details = lastErr.Error()
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "all backends failed",
			"details": details,
		})
		p.m.Error.IncServerErrs()
	}
}
