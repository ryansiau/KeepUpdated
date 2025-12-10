package source

import (
	"fmt"
	"slices"

	"github.com/mitchellh/mapstructure"

	generic_rss "github.com/ryansiau/utilities/go/source/generic-rss"
	"github.com/ryansiau/utilities/go/source/reddit"
	"github.com/ryansiau/utilities/go/source/youtube"
)

// SourceConfig represents the configuration for a content source
type SourceConfig struct {
	Type   string        `yaml:"type"`
	Name   string        `yaml:"name"`
	Config CrawlerConfig `yaml:"config"`
}

func (c *SourceConfig) Validate() error {
	// TODO improve this handling
	validTypes := []string{
		"youtube",
		"reddit",
		"rss",
	}

	// make sure the source type is valid
	if !slices.Contains(validTypes, c.Type) {
		return fmt.Errorf("invalid source type: %s", c.Type)
	}

	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if err := c.Config.Validate(); err != nil {
		return err
	}

	return nil
}

// UnmarshalYAML provides custom unmarshalling for SourceConfig to support interface Config
func (c *SourceConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// temporary structure to capture raw config
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
	// Instantiate concrete source config based on Type
	switch c.Type {
	case "reddit":
		var cfg reddit.RedditCrawlerConfig
		err := mapstructure.Decode(raw.Config, &cfg)
		if err != nil {
			return err
		}
		c.Config = &cfg
	case "youtube":
		var cfg youtube.YouTubeCrawlerConfig
		err := mapstructure.Decode(raw.Config, &cfg)
		if err != nil {
			return err
		}
		c.Config = &cfg
	case "rss":
		var cfg generic_rss.Config
		err := mapstructure.Decode(raw.Config, &cfg)
		if err != nil {
			return err
		}
		c.Config = &cfg
	default:
		return fmt.Errorf("unrecognized config type: %s", c.Type)
	}
	return nil
}

// CrawlerConfig defines the interface for source configurations
type CrawlerConfig interface {
	Validate() error
	IsCrawler()
}
