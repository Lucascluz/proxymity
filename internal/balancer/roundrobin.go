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

func (rr *RoundRobin) NextBackend() (*backend.Backend, error) {
	backends, err := rr.pool.GetHealthyBackends()
	if err != nil {
		return nil, err
	}

	next := atomic.AddUint64(&rr.current, 1)
	idx := int(next-1) % len(backends)
	return backends[idx], err
}
