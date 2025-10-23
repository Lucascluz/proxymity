package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling %s", path)
	}

	// Apply defaults and collect warnings
	warnings := ApplyAllDefaults(&cfg)
	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "CONFIG WARNING: %s\n", w)
	}

	// Validate configs (fatal errors only)
	err = validateBackendConfig(cfg.Backed)
	if err != nil {
		return nil, err
	}

	err = validateProxyConfig(cfg.Proxy)
	if err != nil {
		return nil, err
	}

	err = validateLoadBalancerConfig(cfg.LoadBalancer)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
