package loadbalancer

import (
	"log"
	"proxymity/internal/backend"
	"proxymity/internal/metrics"
)

type LoadBalancer interface {
	NextBackend() (*backend.Backend, error)
	CountAvailableBackends() int
}

func ResolveMethod(method string, pool *backend.Pool, m *metrics.Metrics) LoadBalancer {
	switch method {
	case "round-robin":
		return NewRoundRobin(pool)

	case "random":
		return NewRandom(pool)

	case "least-connections":
		return NewLeastConnections(pool)

	case "weighted":
		return NewWeighted(pool)

	default:
		log.Printf("Error resolving balancer method. Defaulting to round-robin")
		return NewRoundRobin(pool)
	}
}
