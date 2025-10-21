package loadbalancer

import (
	"math/rand"
	"proxymity/internal/backend"
)

type LeastConnections struct {
	pool     *backend.Pool
	connsMap map[*backend.Backend]int
}

func NewLeastConnections(pool *backend.Pool) *LeastConnections {
	return &LeastConnections{
		pool: pool,
	}
}

func (r *LeastConnections) NextBackend() *backend.Backend {
	backends := r.pool.GetHealthyBackends()
	if len(backends) == 0 {
		return nil
	}

	idx := rand.Intn(len(backends))
	return backends[idx]
}
