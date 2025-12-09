package crawler

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	generic_rss "github.com/ryansiau/utilities/go/crawler/generic-rss"
	"github.com/ryansiau/utilities/go/crawler/reddit"
	"github.com/ryansiau/utilities/go/crawler/youtube"
)

// SourceConfig represents the configuration for a content source
type SourceConfig struct {
	Type   string        `yaml:"type"`
	Name   string        `yaml:"name"`
	Config CrawlerConfig `yaml:"config"`
}

// UnmarshalYAML provides custom unmarshalling for SourceConfig to support interface Config
func (s *SourceConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// temporary structure to capture raw config
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
	// Instantiate concrete crawler config based on Type
	switch s.Type {
	case "reddit":
		var cfg reddit.RedditCrawlerConfig
		err := mapstructure.Decode(raw.Config, &cfg)
		if err != nil {
			return err
		}
		s.Config = &cfg
	case "youtube":
		var cfg youtube.YouTubeCrawlerConfig
		err := mapstructure.Decode(raw.Config, &cfg)
		if err != nil {
			return err
		}
		s.Config = &cfg
	case "rss":
		var cfg generic_rss.Config
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

// CrawlerConfig defines the interface for crawler configurations
type CrawlerConfig interface {
	Validate() error
	IsCrawler()
}
