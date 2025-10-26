package metrics

import (
	"sync"
)

type LatencyMetrics struct {
	Min   float64
	Max   float64
	Avg   float64
	P50   float64
	P95   float64
	P99   float64
	mu    sync.Mutex
	count int64 // for calculating average
	sum   float64
}

// RecordLatency records a latency measurement and updates min/max/avg
func (l *LatencyMetrics) RecordLatency(latency float64) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Update min/max
	if latency < l.Min || l.count == 0 {
		l.Min = latency
	}
	if latency > l.Max {
		l.Max = latency
	}

	// Update average
	l.count++
	l.sum += latency
	l.Avg = l.sum / float64(l.count)
}

// UpdatePercentiles updates P50, P95, P99 (placeholder - requires sample collection)
// For now, this is a stub; implement with a sliding window of samples for accuracy
func (l *LatencyMetrics) UpdatePercentiles() {
	l.mu.Lock()
	defer l.mu.Unlock()
	// TODO: Implement percentile calculation from a sample window
	// Example: Collect last N samples, sort, and calculate percentiles
	l.P50 = l.Avg // Placeholder
	l.P95 = l.Avg // Placeholder
	l.P99 = l.Avg // Placeholder
}
