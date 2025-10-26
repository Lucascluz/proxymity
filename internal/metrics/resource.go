package metrics

import "sync"

type ResourceMetrics struct {
	cpuUse float64
	memUSe float64
	mu     sync.Mutex
}

// SetCPUUsage sets the current CPU usage percentage
func (r *ResourceMetrics) SetCPUUsage(usage float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cpuUse = usage
}

// SetMemoryUsage sets the current memory usage percentage
func (r *ResourceMetrics) SetMemoryUsage(usage float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.memUSe = usage
}
