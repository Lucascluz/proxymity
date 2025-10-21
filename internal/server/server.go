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
	"proxymity/internal/proxy"

	"github.com/gin-gonic/gin"
)

type Server struct {
	proxy         *http.Server
	admin         *http.Server
	pool          *backend.Pool
	healthChecker *health.HealthCheck
}

// Create a new http server to receive requests and proxy the to the registered backends.
func New(cfg *config.Config) *Server {

	// Create backend pool
	pool := backend.NewPool()
	for _, bcfg := range cfg.Backed {
		parsedURL, err := url.Parse(bcfg.URL)
		if err != nil {
			// skip invalid backend URL
			continue
		}
		b := &backend.Backend{
			Name: bcfg.Name,
			URL:  parsedURL,
		}
		b.SetAlive(true)
		pool.AddBackend(b)
	}

	// Setup health checker
	hc := health.NewHealthCheck(cfg.HealthCheck, pool)

	// Setup load balancer
	lb := loadbalancer.ResolveMethod(cfg.LoadBalancer.Method, pool)

	// Setup proxy
	p := proxy.NewProxy(lb)

	// Setup proxy router
	proxyRouter := gin.Default()
	proxyRouter.Any("/*path", p.Handler())

	// Setup admin router for proxy information
	adminRouter := gin.New()
	adminRouter.Use(gin.Recovery())
	adminRouter.GET("/health", HealthCheckHandler)
	adminRouter.GET("/status", StatusCheckHandler(pool))
	adminRouter.GET("/config", ConfigCheckHandler(cfg))

	return &Server{
		pool: pool,
		proxy: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.Proxy.Host, cfg.Proxy.Port),
			Handler: proxyRouter,
		},
		admin: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.Proxy.Host, getAdminPort(cfg.Proxy.AdminPort)),
			Handler: adminRouter,
		},
		healthChecker: hc,
	}
}

// getAdminPort returns the admin port or default to 9090
func getAdminPort(port string) string {
	if port == "" {
		return "9090"
	}
	return port
}

func (s *Server) Start() error {
	// Start admin server in background
	go func() {
		log.Printf("Starting admin server on %s", s.admin.Addr)
		if err := s.admin.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Admin server error: %v", err)
		}
	}()

	// Start health checker
	go s.healthChecker.Start()

	// Start proxy server (blocking)
	log.Printf("Starting proxy server on %s", s.proxy.Addr)
	return s.proxy.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {

	// Stoping health checker
	s.healthChecker.Stop()

	// Shutdown both servers
	if err := s.admin.Shutdown(ctx); err != nil {
		log.Printf("Admin server shutdown error: %v", err)
	}

	return s.proxy.Shutdown(ctx)
}
