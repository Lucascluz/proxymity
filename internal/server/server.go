package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"proxymity/internal/backend"
	"proxymity/internal/config"
	"proxymity/internal/loadbalancer"
	"proxymity/internal/proxy"

	"github.com/gin-gonic/gin"
)

type Server struct {
	proxyRouter *gin.Engine
	adminRouter *gin.Engine
	proxyServer *http.Server
	adminServer *http.Server
	pool        *backend.Pool
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
		// Set backend as alive initially
		b.SetAlive(true)
		pool.AddBackend(b)
	}

	// TODO: Handle the error that might come from ResolveMethod
	// Setup load balancer
	lb := loadbalancer.ResolveMethod(cfg.LoadBalancer.Method, pool)

	// Setup proxy
	p := proxy.NewProxy(lb)

	// Setup proxy router (for forwarding traffic to backends)
	proxyRouter := gin.Default()
	proxyRouter.Any("/*path", p.Handler())

	// Setup admin router (for health checks and status)
	adminRouter := gin.New() // Use gin.New() to avoid default middleware
	adminRouter.Use(gin.Recovery())
	adminRouter.GET("/health", HealthCheckHandler)
	adminRouter.GET("/status", StatusCheckHandler(pool))

	return &Server{
		proxyRouter: proxyRouter,
		adminRouter: adminRouter,
		pool:        pool,
		proxyServer: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.Proxy.Host, cfg.Proxy.Port),
			Handler: proxyRouter,
		},
		adminServer: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.Proxy.Host, getAdminPort(cfg.Proxy.AdminPort)),
			Handler: adminRouter,
		},
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
		log.Printf("Starting admin server on %s", s.adminServer.Addr)
		if err := s.adminServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Admin server error: %v", err)
		}
	}()

	// Start proxy server (blocking)
	log.Printf("Starting proxy server on %s", s.proxyServer.Addr)
	return s.proxyServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	// Shutdown both servers
	if err := s.adminServer.Shutdown(ctx); err != nil {
		log.Printf("Admin server shutdown error: %v", err)
	}
	return s.proxyServer.Shutdown(ctx)
}
