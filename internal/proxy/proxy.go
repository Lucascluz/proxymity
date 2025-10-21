package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"proxymity/internal/loadbalancer"

	"github.com/gin-gonic/gin"
)

type Proxy struct {
	lb loadbalancer.LoadBalancer
}

func NewProxy(lb loadbalancer.LoadBalancer) *Proxy {
	return &Proxy{lb: lb}
}

func (p *Proxy) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		backend := p.lb.NextBackend()
		if backend == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "no healthy backends available",
			})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(backend.URL)

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Backend %s error: %v", backend.Name, err)
			backend.SetAlive(false)

			c.JSON(http.StatusBadGateway, gin.H{
				"error": "backend unavailable",
			})
		}

		proxy.ServeHTTP(c.Writer, c.Request)
		backend.AddConnection()
	}
}
