package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"proxymity/internal/backend"
	"proxymity/internal/config"
	"proxymity/internal/loadbalancer"
	"proxymity/internal/proxy"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router     *gin.Engine
	httpServer *http.Server
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
	lb, _ := loadbalancer.ResolveMethod(cfg.LoadBalancer.Method, pool)

	// Setup proxy
	p := proxy.NewProxy(lb)

	// Setup router
	router := gin.Default()

	router.Any("/*path", p.Handler())

	return &Server{
		router: router,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.Proxy.Host, cfg.Proxy.Port),
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(context.Context) {

}
