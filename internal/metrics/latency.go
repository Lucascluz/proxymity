package metrics

import (
	"sync"
)

type LatencyMetrics struct {
	Min   float64
	Max   float64
	Avg   float64
	
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

