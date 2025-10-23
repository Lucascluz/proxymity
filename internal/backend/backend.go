package backend

import (
	"net/url"
	"sync"
)

type Backend struct {

	// Arbitrary name for the backend
	Name string

	// Root path of the backend
	Host *url.URL

	// Path to the health check endpoint with /. Default to /health. If root path, insert /
	Health string

	// Define if will receive requests regardless of being alive or not
	Enabled bool

	alive  bool
	conns  int
	weight int
	mu     sync.Mutex
}

// Returns the current live state of the backendU+
func (b *Backend) IsAlive() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.alive
}

// Updates the live state of the backend
func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.alive = alive
}

// Return the current number of connections the backend has stabilished in his lifetime
func (b *Backend) GetConnections() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.conns
}

// Increases the connection counter
func (b *Backend) AddConnection() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.conns++
}

// Returns the connection bias wieght from that specific backend
func (b *Backend) GetWeight() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.weight
}
