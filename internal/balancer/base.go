package loadbalancer

import "proxymity/internal/backend"

// BaseLoadBalancer provides shared logic for all load balancers
type BaseLoadBalancer struct {
	pool *backend.Pool
}

func (b *BaseLoadBalancer) CountAvailableBackends() int {
	backends, err := b.pool.GetAvailableBackends()
	if err != nil {
		return 0
	}
	return len(backends)
}
