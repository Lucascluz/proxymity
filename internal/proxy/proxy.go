package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"proxymity/internal/balancer"
	"proxymity/internal/metrics"

	"github.com/gin-gonic/gin"
)

type Proxy struct {
	lb loadbalancer.LoadBalancer
	m *metrics.Metrics
}

func NewProxy(lb loadbalancer.LoadBalancer, m *metrics.Metrics) *Proxy {
	return &Proxy{lb: lb}
}

func (p *Proxy) Proxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		backend, err := p.lb.NextBackend()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": err.Error(),
			})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(backend.Host)

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Errora at %s: %v", backend.Name, err)
			backend.SetAlive(false)

			c.JSON(http.StatusBadGateway, gin.H{
				"error": "backend unavailable",
			})
		}

		proxy.ServeHTTP(c.Writer, c.Request)
		backend.AddConnection()
	}
}
