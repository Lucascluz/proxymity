package metrics

import "sync/atomic"

type TrafficMetrics struct {
	Requests  int64
	Responses int64
	BytesIn   int64
	BytesOut  int64
}

// IncRequests increments the total requests counter
func (t *TrafficMetrics) IncRequests() {
	atomic.AddInt64(&t.Requests, 1)
}

// IncResponses increments the total responses counter
func (t *TrafficMetrics) IncResponses() {
	atomic.AddInt64(&t.Responses, 1)
}

// AddBytesIn adds to the bytes received counter
func (t *TrafficMetrics) AddBytesIn(bytes int64) {
	atomic.AddInt64(&t.BytesIn, bytes)
}

// AddBytesOut adds to the bytes sent counter
func (t *TrafficMetrics) AddBytesOut(bytes int64) {
	atomic.AddInt64(&t.BytesOut, bytes)
}
