package loadbalancer

import (
	"proxymity/internal/backend"
	"sync/atomic"
)

type RoundRobin struct {
	pool    *backend.Pool
	current uint64
}

func NewRoundRobin(pool *backend.Pool) *RoundRobin {
	return &RoundRobin{
		pool: pool,
	}
}

func (rr *RoundRobin) NextBackend() *backend.Backend {
	backends := rr.pool.GetHealthyBackends()
	if len(backends) == 0 {
		return nil
	}

	next := atomic.AddUint64(&rr.current, 1)
	idx := int(next-1) % len(backends)
	return backends[idx]
}
