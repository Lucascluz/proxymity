package metrics

import (
	"sync/atomic"
)

type ErrorMetrics struct {
	Total    int64
	Client   int64
	Proxy    int64
	Server   int64
	Timeouts int64
	Retrys   int64
}

// IncClientErrs increments the client errors counter
func (e *ErrorMetrics) IncClientErrs() {
	atomic.AddInt64(&e.Client, 1)
	atomic.AddInt64(&e.Total, 1)
}

// IncServerErrs increments the server errors counter
func (e *ErrorMetrics) IncServerErrs() {
	atomic.AddInt64(&e.Server, 1)
	atomic.AddInt64(&e.Total, 1)
}

// IncTimeoutsErrs increments the timeouts counter
func (e *ErrorMetrics) IncTimeoutsErrs() {
	atomic.AddInt64(&e.Timeouts, 1)
	atomic.AddInt64(&e.Total, 1)
}

// IncProxyErrs increments the proxy errors counter
func (e *ErrorMetrics) IncProxyErrs() {
	atomic.AddInt64(&e.Proxy, 1)
	atomic.AddInt64(&e.Total, 1)
}

// IncRetrysErrs increments the retrys counter
func (e *ErrorMetrics) IncRetrysErrs() {
	atomic.AddInt64(&e.Retrys, 1)
}
