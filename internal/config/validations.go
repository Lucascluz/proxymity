package config

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

func validateBackendConfig(cfg []BackendConfig) error {

	// At least one backend
	if len(cfg) < 1 {
		return errors.New("no backends conigured")
	}

	for _, b := range cfg {
		// Backend name (non-empty)
		if b.Name == "" {
			return errors.New("all backends should have names attributed")
		}

		// Backend URL (non-empty, valid format)
		if !isValidUrl(b.Host) {
			return fmt.Errorf("invalid host url for %s", b.Name)
		}

		// Backend health check path (starts with /, valid format, if empty default to /health)
		if b.Health != "/" {

		}
	}

	return nil
}

func validateProxyConfig(cfg ProxyConfig) error {

	// Backend URL (non-empty, valid format)
	if !isValidUrl(cfg.Host) {
		return fmt.Errorf("%s is not a valid host", cfg.Host)
	}

	// Validate proxy port
	err := isValidPort(cfg.Port)
	if err != nil {
		return fmt.Errorf("invalid proxy port")
	}

	// Validate admin port
	err = isValidPort(cfg.AdminPort)
	if err != nil {
		return fmt.Errorf("invalid proxy admin port")
	}

	return nil
}

func validateHealthCheckerConfig(cfg HealthCheckConfig) {

}

func isValidUrl(str string) bool {

	if str == "0.0.0.0" || str == "localhost" {
		return true
	}

	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func isValidPort(port string) error {

	// Convert port to number
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("error converting port %s", port)
	}

	// Validate port value
	if portNum < 0 || portNum > 65535 {
		return fmt.Errorf("invalid proxy port")
	}

	return nil
}
