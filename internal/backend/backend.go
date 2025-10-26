package backend

import (
	"fmt"
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

	// Whether the backend is enabled to receive traffic
	Enabled bool

	// Weight for weighted load balancing methods
	Weight int

	healthy bool
	conns   int
	mu      sync.Mutex
}

// Returns the current live state of the backendU+
func (b *Backend) IsHealthy() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.healthy
}

// Updates the live state of the backend
func (b *Backend) SetHealthy(Healthy bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.healthy = Healthy
}

func (b *Backend) IsEnabled() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Enabled
}

// Enables or disables the backend from receiving traffic
func (b *Backend) SetEnabled(Enabled bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Enabled = Enabled
}

func (b *Backend) IsAvailable() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.Enabled {
		fmt.Printf("Backend %s is disabled\n", b.Name)
		return false
	}

	if !b.healthy {
		fmt.Printf("Backend %s is unHealthy\n", b.Name)
		return false
	}

	return b.Enabled && b.healthy
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
	return b.Weight
}
