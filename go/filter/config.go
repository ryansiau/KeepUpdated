package filter

import (
	"fmt"
)

type FilterConfig struct {
	Name      string      `yaml:"name"`
	Type      string      `yaml:"type"`
	Config    interface{} `yaml:"config,omitempty"`
	Sources   []string    `yaml:"sources"`
	Notifiers []string    `yaml:"notifiers"`
}

// Validate validates the filter configuration
func (f *FilterConfig) Validate() error {
	if f.Name == "" {
		return fmt.Errorf("filter name is required")
	}
	if f.Type == "" {
		return fmt.Errorf("filter type is required")
	}
	return nil
}
