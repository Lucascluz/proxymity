package backend

import (
	"fmt"
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

func (p *Pool) GetAvailableBackends() ([]*Backend, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	available := make([]*Backend, 0)
	for _, b := range p.backends {
		if b.IsAvailable() {
			available = append(available, b)
			fmt.Printf("Backend %s is available\n", b.Name)
		} else {
			fmt.Printf("Backend %s is not available\n", b.Name)
		}
	}

	if len(available) == 0 {
		return nil, fmt.Errorf("no backends available")
	}

	return available, nil
}
