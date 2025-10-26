package metrics

import (
	"sync/atomic"
)

type ErrorMetrics struct {
	total    int64
	client   int64
	proxy    int64
	server   int64
	timeouts int64
	retrys   int64
}

// Incclient increments the client errors counter
func (e *ErrorMetrics) IncClientErrs() {
	atomic.AddInt64(&e.client, 1)
	atomic.AddInt64(&e.total, 1)
}

// IncServer increments the server errors counter
func (e *ErrorMetrics) IncServerErrs() {
	atomic.AddInt64(&e.server, 1)
	atomic.AddInt64(&e.total, 1)
}

// IncTimeouts increments the timeouts counter
func (e *ErrorMetrics) IncTimeoutsErrs() {
	atomic.AddInt64(&e.timeouts, 1)
	atomic.AddInt64(&e.total, 1)
}
