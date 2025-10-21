package loadbalancer

import (
	"math/rand"
	"proxymity/internal/backend"
)

type Random struct {
	pool *backend.Pool
}

func NewRandom(pool *backend.Pool) *Random {
	return &Random{
		pool: pool,
	}
}

func (r *Random) NextBackend() (*backend.Backend, error) {
	backends, err := r.pool.GetHealthyBackends()
	if err != nil {
		return nil, err
	}

	idx := rand.Intn(len(backends))
	return backends[idx], err
}
