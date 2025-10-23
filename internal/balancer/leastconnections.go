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

func (lc *LeastConnections) NextBackend() (*backend.Backend, error) {
	backends, err := lc.pool.GetAvailableBackends()
	if err != nil {
		return nil, err
	}

	idx := rand.Intn(len(backends))
	return backends[idx], nil
}
