package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	loadbalancer "proxymity/internal/balancer"
	"proxymity/internal/metrics"
	"strings"
	"time"

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

// responseWriter wraps gin.ResponseWriter to track bytes written and status
type responseWriter struct {
	gin.ResponseWriter
	bytesWritten int64
	status       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK // Default if not set
	}
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
				// Classify proxy errors
				if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline") {
					p.m.Error.IncTimeoutsErrs()
				} else {
					p.m.Error.IncProxyErrs()
				}
			}

			// Wrap the response writer to track bytes out and status
			rw := &responseWriter{ResponseWriter: c.Writer}

			latency := time.Now()
			proxy.ServeHTTP(rw, c.Request)
			backend.AddConnection()

			// Update metrics
			p.m.Traffic.IncRequests()
			p.m.Traffic.AddBytesIn(c.Request.ContentLength)
			p.m.Traffic.AddBytesOut(rw.bytesWritten)

			// Classify based on response status
			if rw.status >= 400 && rw.status < 500 {
				p.m.Error.IncClientErrs()
			} else if rw.status >= 500 {
				p.m.Error.IncServerErrs()
			}

			// If no error was set by ErrorHandler, request succeeded
			if lastErr == nil {
				p.m.Latency.RecordLatency(float64(time.Since(latency)))
				p.m.Traffic.IncResponses()
				return
			}

			tried++
			p.m.Error.IncRetrysErrs()
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
		p.m.Error.IncProxyErrs()
	}
}
