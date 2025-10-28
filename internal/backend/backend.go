package backend

import (
	"net/url"
	"sync"
	"time"
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
	checked time.Time
	backoff time.Duration
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
		return false
	}

	if !b.healthy {
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

func (b *Backend) GetChecked() time.Time {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.checked
}

func (b *Backend) SetChecked() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.checked = time.Now()
}

func (b *Backend) GetBackoff() time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.backoff
}

func (b *Backend) ExpBackof() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.backoff == 0 {
		b.backoff = 1 * time.Second
		return
	}

	if b.backoff >= 1*time.Minute {
		return
	}

	b.backoff = b.backoff * 2
}

func (b *Backend) ResetBackof() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.backoff = 0
}
