package filter

import (
	"fmt"
	"slices"

	"github.com/mitchellh/mapstructure"

	"github.com/ryansiau/utilities/go/filter/metadata"
	"github.com/ryansiau/utilities/go/filter/title"
)

type BaseConfig struct {
	Name   string       `yaml:"name"`
	Type   string       `yaml:"type"`
	Config FilterConfig `yaml:"config,omitempty"`
}

// Validate validates the filter configuration
func (c *BaseConfig) Validate() error {
	validTypes := []string{
		"metadata",
		"title",
	}
	if !slices.Contains(validTypes, c.Type) {
		return fmt.Errorf("invalid filter type: %s", c.Type)
	}

	if c.Name == "" {
		return fmt.Errorf("filter name is required")
	}
	if err := c.Config.Validate(); err != nil {
		return err
	}
	return nil
}

func (s *BaseConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw struct {
		Type   string                 `yaml:"type"`
		Name   string                 `yaml:"name"`
		Config map[string]interface{} `yaml:"config"`
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}
	s.Type = raw.Type
	s.Name = raw.Name

	switch s.Type {
	case "title":
		var cfg title.Config
		err := mapstructure.Decode(raw.Config, &cfg)
		if err != nil {
			return err
		}
		s.Config = &cfg
	case "metadata":
		var cfg metadata.Config
		err := mapstructure.Decode(raw.Config, &cfg)
		if err != nil {
			return err
		}

		s.Config = &cfg
	default:
		return fmt.Errorf("unrecognized config type: %s", s.Type)
	}
	return nil
}

type FilterConfig interface {
	Validate() error
	IsFilterConfig()
}
