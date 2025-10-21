package backend

import (
	"net/url"
	"sync"
)

type Backend struct {
	Name   string
	URL    *url.URL
	alive  bool
	conns  int
	weight int
	mu     sync.Mutex
}

func (b *Backend) IsAlive() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.alive
}

func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.alive = alive
}

func (b *Backend) GetConnections() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.conns
}

func (b *Backend) AddConnection() {
	b.conns++
}

func (b *Backend) GetWeight() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.weight
}
