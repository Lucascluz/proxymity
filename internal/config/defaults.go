package config

import "fmt"

// Default values
const (
	DefaultLoadBalancerMethod = "round-robin"
	DefaultHealthCheckPath    = "/health"
	DefaultBackendWeight      = 1
	DefaultHealthInterval     = 30 // seconds
	DefaultHealthTimeout      = 5  // seconds
)

// Applies default values to backend configurations and returns a slice of warning messages for any defaults that were applied
func ApplyBackendDefaults(backends []BackendConfig) []string {
	warnings := []string{}

	for _, b := range backends {
		// Default health check path
		if b.Health == "" {
			b.Health = DefaultHealthCheckPath
			warnings = append(warnings, fmt.Sprintf("Backend '%s': health check path not specified, using default: %s", b.Name, DefaultHealthCheckPath))
		}

		// Default weight
		if b.Weight <= 0 {
			oldWeight := b.Weight
			b.Weight = DefaultBackendWeight
			if oldWeight != 0 {
				warnings = append(warnings, fmt.Sprintf("Backend '%s': invalid weight %d, using default: %d", b.Name, oldWeight, DefaultBackendWeight))
			}
		}

	}

	return warnings
}

// Applies default values to load balancer configuration and returns a slice of warning messages for any defaults that were applied
func ApplyLoadBalancerDefaults(lb *LoadBalancerConfig) []string {
	warnings := []string{}

	validMethods := map[string]bool{
		"round-robin":       true,
		"least-connections": true,
		"weighted":          true,
		"random":            true,
	}

	// Default load balancer method
	if lb.Method == "" {
		lb.Method = DefaultLoadBalancerMethod
		warnings = append(warnings, fmt.Sprintf("Load balancer method not specified, using default: %s", DefaultLoadBalancerMethod))
	} else if !validMethods[lb.Method] {
		oldMethod := lb.Method
		lb.Method = DefaultLoadBalancerMethod
		warnings = append(warnings, fmt.Sprintf("Invalid load balancer method '%s', using default: %s", oldMethod, DefaultLoadBalancerMethod))
	}

	return warnings
}

// Applies default values to health check configuration and returns a slice of warning messages for any defaults that were applied
func ApplyHealthCheckDefaults(hc *HealthCheckConfig) []string {
	warnings := []string{}

	// Default interval
	if hc.Interval == 0 {
		hc.Interval = DefaultHealthInterval
		warnings = append(warnings, fmt.Sprintf("Health check interval not specified, using default: %d seconds", DefaultHealthInterval))
	}

	// Default timeout
	if hc.TimeOut == 0 {
		hc.TimeOut = DefaultHealthTimeout
		warnings = append(warnings, fmt.Sprintf("Health check timeout not specified, using default: %d seconds", DefaultHealthTimeout))
	}

	// Validate that timeout is less than interval
	if hc.TimeOut >= hc.Interval {
		warnings = append(warnings, fmt.Sprintf("Warning: Health check timeout (%d) should be less than interval (%d)", hc.TimeOut, hc.Interval))
	}

	return warnings
}

// Applies all default values to the configuration and returns a slice of all warning messages
func ApplyAllDefaults(cfg *Config) []string {
	warnings := []string{}

	warnings = append(warnings, ApplyBackendDefaults(cfg.Backed)...)
	warnings = append(warnings, ApplyLoadBalancerDefaults(&cfg.LoadBalancer)...)
	warnings = append(warnings, ApplyHealthCheckDefaults(&cfg.HealthCheck)...)

	return warnings
}
