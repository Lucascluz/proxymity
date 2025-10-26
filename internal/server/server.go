package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"proxymity/internal/backend"
	loadbalancer "proxymity/internal/balancer"
	"proxymity/internal/config"
	"proxymity/internal/health"
	"proxymity/internal/metrics"
	"proxymity/internal/proxy"

	"github.com/gin-gonic/gin"
)

type Server struct {
	proxy         *http.Server
	pool          *backend.Pool
	healthChecker *health.HealthChecker
	metrics       *metrics.Metrics
}

// Create a new http server to receive requests and proxy the to the registered backends.
func New(cfg *config.Config) *Server {

	// Setup metrics
	m := metrics.NewMetrics()

	// Create backend pool
	pool := backend.NewPool()
	for _, bcfg := range cfg.Backend {
		parsedURL, err := url.Parse(bcfg.Host)
		if err != nil {
			// skip invalid backend URL
			continue
		}
		b := &backend.Backend{
			Name:    bcfg.Name,
			Host:    parsedURL,
			Health:  bcfg.Health,
			Enabled: *bcfg.Enabled,
			Weight:  bcfg.Weight,
		}
		b.SetHealthy(true)
		pool.AddBackend(b)
	}

	// Setup health checker
	hc := health.NewHealthChecker(cfg.HealthCheck, pool, m)

	// Setup load balancer
	lb := loadbalancer.ResolveMethod(cfg.LoadBalancer.Method, pool, m)

	// Setup proxy
	p := proxy.NewProxy(lb, m)

	// Setup proxy router
	pRouter := gin.Default()
	pRouter.GET("/api/proxy/health", Health)
	pRouter.GET("/api/proxy/status", Status(pool))
	pRouter.GET("/api/proxy/config", Config(cfg))
	pRouter.GET("/metrics", Metrics(m))
	pRouter.NoRoute(p.Proxy())

	return &Server{
		pool: pool,
		proxy: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.Proxy.Host, cfg.Proxy.Port),
			Handler: pRouter,
		},
		healthChecker: hc,
		metrics:       m,
	}
}

func (s *Server) Start() error {

	// Start health checker
	go s.healthChecker.Start()

	// Start proxy server (blocking)
	log.Printf("Starting proxy server on %s", s.proxy.Addr)
	return s.proxy.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {

	// Stoping health checker
	s.healthChecker.Stop()

	return s.proxy.Shutdown(ctx)
}
