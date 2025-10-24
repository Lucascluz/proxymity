package config

import "proxymity/internal/metrics"

type Config struct {
	Proxy        ProxyConfig        `yaml:"proxy"`
	Backed       []BackendConfig    `yaml:"backend"`
	LoadBalancer LoadBalancerConfig `yaml:"load-balancer"`
	HealthCheck  HealthCheckConfig  `yaml:"health-check"`

	m *metrics.Metrics
}

type ProxyConfig struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	AdminPort string `yaml:"admin_port"`
}

type BackendConfig struct {
	Name    string `yaml:"name"`
	Host    string `yaml:"url"`
	Health  string `yaml:"health"`
	Weight  int    `yaml:"weight"`
	Enabled bool   `yaml:"enabled"`
}

type LoadBalancerConfig struct {
	Method string `yaml:"method"`
}

type HealthCheckConfig struct {
	Interval uint // Interval between checks in seconds
	TimeOut  uint // Health check timeout in seconds
}
