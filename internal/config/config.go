package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Proxy        ProxyConfig        `yaml:"proxy"`
	Backed       []BackendConfig    `yaml:"backend"`
	LoadBalancer LoadBalancerConfig `yaml:"load-balancer"`
}

type ProxyConfig struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	AdminPort string `yaml:"admin_port"`
}

type BackendConfig struct {
	Name    string `yaml:"name"`
	URL     string `yaml:"url"`
	Weight  int    `yaml:"weight"`
	Enabled bool   `yaml:"enabled"`
}

type LoadBalancerConfig struct {
	Method string `yaml:"method"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}
