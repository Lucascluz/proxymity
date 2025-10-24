package metrics

type Metrics struct {
	Traffic      *TrafficMetrics
	Latency      *LatencyMetrics
	Error        *ErrorMetrics
	Backend      *BackendMetrics
	LoadBalancer *LoadBalancerMetrics
	Resource     *ResourceMetrics
}

func NewMetrics() *Metrics {
	return &Metrics{
		Traffic:      &TrafficMetrics{},
		Latency:      &LatencyMetrics{},
		Error:        &ErrorMetrics{},
		Backend:      &BackendMetrics{},
		LoadBalancer: &LoadBalancerMetrics{RequestsPerBackend: make(map[string]int64)},
		Resource:     &ResourceMetrics{},
	}
}