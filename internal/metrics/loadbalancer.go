package metrics

type LoadBalancerMetrics struct {
	RequestsPerBackend map[string]int64
	ActiveConnections  int
}
