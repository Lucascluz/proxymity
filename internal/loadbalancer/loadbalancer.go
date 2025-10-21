package loadbalancer

import (
	"log"
	"proxymity/internal/backend"
)

type LoadBalancer interface {
	NextBackend() (*backend.Backend, error)
}

func ResolveMethod(method string, pool *backend.Pool) LoadBalancer {
	switch method {
	case "round-robin":
		st := NewRoundRobin(pool)
		return st

	case "random":
		st := NewRandom(pool)
		return st
	default:
		log.Printf("No load balancing method recognized, defaulting to round-robin")
		return NewRoundRobin(pool)
	}
}
