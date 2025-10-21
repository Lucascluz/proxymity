package backend

import (
	"errors"
	"sync"
)

type Pool struct {
	backends []*Backend
	mu       sync.Mutex
}

func NewPool() *Pool {
	return &Pool{
		backends: make([]*Backend, 0),
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
		return nil, errors.New("no healthy backends available")
	}
	return healthy, nil
}
