package notification

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/ryansiau/utilities/go/model"
	"github.com/ryansiau/utilities/go/notification/discord"
	"github.com/ryansiau/utilities/go/notification/ntfy"
)

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

func (c *BaseConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw struct {
		Type   string                 `yaml:"type"`
		Name   string                 `yaml:"name"`
		Config map[string]interface{} `yaml:"config"`
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}
	c.Type = raw.Type
	c.Name = raw.Name

	switch c.Type {
	case "discord":
		var cfg discord.Config
		err := mapstructure.Decode(raw.Config, &cfg)
		if err != nil {
			return err
		}
		c.Config = &cfg

	case "ntfy":
		var cfg ntfy.Config
		err := mapstructure.Decode(raw.Config, &cfg)
		if err != nil {
			return err
		}
		c.Config = &cfg

	case "default":
		// this is fine enough just to avoid being caught by default and returning an error

	default:
		return fmt.Errorf("unrecognized config type: %s", c.Type)
	}
	return nil
}

type NotifierConfig interface {
	Validate() error
	IsNotifierConfig()
	Build() (model.Notifier, error)
}
