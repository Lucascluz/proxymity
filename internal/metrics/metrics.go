package metrics

type Metrics struct {
	Traffic      *TrafficMetrics
	Latency      *LatencyMetrics
	Error        *ErrorMetrics
	Resource     *ResourceMetrics
}

func NewMetrics() *Metrics {
	return &Metrics{
		Traffic:      &TrafficMetrics{},
		Latency:      &LatencyMetrics{},
		Error:        &ErrorMetrics{},
		Resource:     &ResourceMetrics{},
	}
}