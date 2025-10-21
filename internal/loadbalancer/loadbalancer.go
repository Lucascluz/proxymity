package loadbalancer

import (
	"errors"
	"proxymity/internal/backend"
)

type LoadBalancer interface {
	NextBackend() *backend.Backend
}

func ResolveMethod(method string, pool *backend.Pool) (LoadBalancer, error) {
	switch method {
	case "round-robin":
		st := NewRoundRobin(pool)
		return st, nil

	case "random":
		st := NewRandom(pool)
		return st, nil
	}
	return nil, errors.New("method not recognized")
}
