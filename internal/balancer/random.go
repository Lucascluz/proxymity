package loadbalancer

import (
	"math/rand"
	"proxymity/internal/backend"
)

type Random struct {
	BaseLoadBalancer
}

func NewRandom(pool *backend.Pool) *Random {
	return &Random{
		BaseLoadBalancer: BaseLoadBalancer{pool: pool},
	}
}

func (r *Random) NextBackend() (*backend.Backend, error) {
	backends, err := r.pool.GetAvailableBackends()
	if err != nil {
		return nil, err
	}

	idx := rand.Intn(len(backends))
	return backends[idx], err
}
