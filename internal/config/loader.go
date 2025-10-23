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

	err = validateBackendConfig(cfg.Backed)
	if err != nil {
		return nil, err
	}

	err = validateProxyConfig(cfg.Proxy)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
