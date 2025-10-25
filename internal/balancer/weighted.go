package loadbalancer

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"proxymity/internal/backend"
)

type Weighted struct {
	BaseLoadBalancer
}

func NewWeighted(pool *backend.Pool) *Weighted {
	return &Weighted{
		BaseLoadBalancer: BaseLoadBalancer{pool: pool},
	}
}

func (rr *Weighted) NextBackend() (*backend.Backend, error) {
	backends, err := rr.pool.GetAvailableBackends()
	if err != nil {
		return nil, err
	}

	// Calculate total weight
	total := 0
	for _, b := range backends {
		total += b.GetWeight()
	}
	if total <= 0 {
		return nil, fmt.Errorf("total weight is zero")
	}

	max := big.NewInt(int64(total))
	nBig, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, fmt.Errorf("error generating random weight coefficient: %w", err)
	}
	coeff := int(nBig.Int64()) + 1

	for _, b := range backends {
		coeff -= b.GetWeight()
		if coeff <= 0 {
			return b, nil
		}
	}

	// Fallback: return last backend if any
	if len(backends) > 0 {
		return backends[len(backends)-1], nil
	}

	return nil, fmt.Errorf("no backend selected")
}
