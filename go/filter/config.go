package filter

import (
	"fmt"
)

type BaseConfig struct {
	Name   string       `yaml:"name"`
	Type   string       `yaml:"type"`
	Config FilterConfig `yaml:"config,omitempty"`
}

// Validate validates the filter configuration
func (c *BaseConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("filter name is required")
	}
	if c.Type == "" {
		return fmt.Errorf("filter type is required")
	}
	if err := c.Config.Validate(); err != nil {
		return err
	}
	return nil
}

type FilterConfig interface {
	Validate() error
	IsFilterConfig()
}
