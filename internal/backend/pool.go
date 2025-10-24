package backend

import (
	"fmt"
	"proxymity/internal/metrics"
	"sync"
)

type Pool struct {
	backends []*Backend
	metrics  *metrics.Metrics
	mu       sync.Mutex
}

func NewPool(m *metrics.Metrics) *Pool {
	return &Pool{
		backends: make([]*Backend, 0),
		metrics:  m,
	}
}

func (p *Pool) AddBackend(b *Backend) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.backends = append(p.backends, b)
}

func (p *Pool) GetBackends() []*Backend {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.backends
}

func (p *Pool) GetHealthyBackends() ([]*Backend, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	healthy := make([]*Backend, 0)
	for _, b := range p.backends {
		if b.IsAlive() {
			healthy = append(healthy, b)
		}
	}

	if len(healthy) == 0 {
		return nil, fmt.Errorf("no healthy backends available")
	}
	return healthy, nil
}

func (p *Pool) GetAvailableBackends() ([]*Backend, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	healthy, err := p.GetHealthyBackends()
	if err != nil {
		return nil, err
	}
	available := make([]*Backend, 0)
	for _, b := range healthy {
		if b.IsAlive() {
			available = append(available, b)
		}
	}

	if len(available) == 0 {
		return nil, fmt.Errorf("there are no backends available at the moment, retrying")
	}
	return available, nil
}
