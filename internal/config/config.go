package config

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Host     string
	Timeout  time.Duration
	Raw      bool
	Insecure bool
}

func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("--host must not be empty")
	}
	if !strings.HasPrefix(c.Host, "http://") && !strings.HasPrefix(c.Host, "https://") {
		return fmt.Errorf("--host must start with http:// or https://")
	}
	return nil
}
