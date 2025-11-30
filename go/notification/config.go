package notification

import "errors"

type BaseConfig struct {
	Type   string         `yaml:"type"`
	Name   string         `yaml:"name"`
	Config NotifierConfig `yaml:"config,omitempty"`
}

func (c *BaseConfig) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	if c.Type == "" {
		return errors.New("type is required")
	}
	if err := c.Config.Validate(); err != nil {
		return err
	}
	return nil
}

type NotifierConfig interface {
	Validate() error
	IsNotifierConfig()
}
