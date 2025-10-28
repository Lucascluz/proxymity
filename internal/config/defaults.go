package config

import "fmt"

// Default values
const (
	DefaultProxyHost          = "0.0.0.0"
	DefaultProxyPort          = "8080"
	DefaultAdminPort          = "9090"
	DefaultLoadBalancerMethod = "round-robin"
	DefaultHealthCheckPath    = "/health"
	DefaultBackendWeight      = 1
	DefaultHealthInterval     = 10 // seconds
	DefaultHealthTimeout      = 3  // seconds
)

// Applies default values to backend configurations and returns a slice of warning messages for any defaults that were applied
func ApplyBackendDefaults(backends *[]BackendConfig) []string {
	warnings := []string{}

	for i := range *backends {
		if (*backends)[i].Health == "" {
			(*backends)[i].Health = DefaultHealthCheckPath
			warnings = append(warnings, fmt.Sprintf("Backend '%s': health check path not specified, using default: %s", (*backends)[i].Name, DefaultHealthCheckPath))
		}

		if (*backends)[i].Weight <= 0 {
			oldWeight := (*backends)[i].Weight
			(*backends)[i].Weight = DefaultBackendWeight
			if oldWeight != 0 {
				warnings = append(warnings, fmt.Sprintf("Backend '%s': invalid weight %d, using default: %d", (*backends)[i].Name, oldWeight, DefaultBackendWeight))
			}
		}

		if (*backends)[i].Enabled == nil {
			defaultEnabled := true
			(*backends)[i].Enabled = &defaultEnabled
			warnings = append(warnings, fmt.Sprintf("Backend '%s': enabled not specified, using default: true", (*backends)[i].Name))
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

	if hc.Interval == 0 {
		hc.Interval = DefaultHealthInterval
		warnings = append(warnings, fmt.Sprintf("Health check interval not specified, using default: %d seconds", DefaultHealthInterval))
	}

	if hc.TimeOut == 0 {
		hc.TimeOut = DefaultHealthTimeout
		warnings = append(warnings, fmt.Sprintf("Health check timeout not specified, using default: %d seconds", DefaultHealthTimeout))
	}

	if hc.TimeOut >= hc.Interval {
		warnings = append(warnings, fmt.Sprintf("Warning: Health check timeout (%d) should be less than interval (%d)", hc.TimeOut, hc.Interval))
	}

	return warnings
}

func ApplyProxyDefaults(p *ProxyConfig) []string {
	warnings := []string{}

	if p.Host == "" {
		p.Host = DefaultProxyHost
		warnings = append(warnings, fmt.Sprintf("Proxy host not specified, using default: %s", DefaultProxyHost))
	}

	if p.Port == "" {
		p.Port = DefaultProxyPort
		warnings = append(warnings, fmt.Sprintf("Proxy Port not specified, using default: %s", DefaultProxyPort))
	}

	if p.AdminPort == "" {
		p.AdminPort = DefaultAdminPort
		warnings = append(warnings, fmt.Sprintf("Proxy AdminPort not specified, using default: %s", DefaultAdminPort))
	}

	return warnings
}

// Applies all default values to the configuration and returns a slice of all warning messages
func ApplyAllDefaults(cfg *Config) []string {
	warnings := []string{}

	warnings = append(warnings, ApplyBackendDefaults(&cfg.Backend)...)
	warnings = append(warnings, ApplyLoadBalancerDefaults(&cfg.LoadBalancer)...)
	warnings = append(warnings, ApplyHealthCheckDefaults(&cfg.HealthCheck)...)
	warnings = append(warnings, ApplyProxyDefaults(&cfg.Proxy)...)

	return warnings
}
