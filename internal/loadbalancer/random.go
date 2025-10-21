package loadbalancer

import (
	"math/rand"
	"proxymity/internal/backend"
)

type Random struct {
	pool    *backend.Pool
}

func NewRandom(pool *backend.Pool) *Random {
	return &Random{
		pool: pool,
	}
}

func (r *Random) NextBackend() *backend.Backend {
	backends := r.pool.GetHealthyBackends()
	if len(backends) == 0 {
		return nil
	}

	idx := rand.Intn(len(backends))
	return backends[idx]
}
