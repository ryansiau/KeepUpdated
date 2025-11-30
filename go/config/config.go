package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/ryansiau/utilities/go/crawler"
	"github.com/ryansiau/utilities/go/crawler/reddit"
	"github.com/ryansiau/utilities/go/crawler/youtube"
	"github.com/ryansiau/utilities/go/filter"
	"github.com/ryansiau/utilities/go/notification"
)

// Config represents the entire configuration
type Config struct {
	Sources     []crawler.SourceConfig        `yaml:"sources"`
	Notifiers   []notification.NotifierConfig `yaml:"notifiers"`
	Filters     []filter.BaseConfig           `yaml:"filters"`
	Credentials Credentials                   `yaml:"credentials"`
}

type Credentials struct {
	YoutubeApiKey string `yaml:"youtube_api_key"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(filepath string) (*Config, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var cfg Config
	// Unmarshal the entire file
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	// Validate the loaded config
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	return &cfg, nil
}

// SaveConfig saves configuration to a YAML file
func SaveConfig(cfg *Config, filepath string) error {
	content, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := os.WriteFile(filepath, content, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

// Validate validates the entire configuration structure
func (c *Config) Validate() error {
	if err := c.ValidateSources(); err != nil {
		return err
	}
	if err := c.ValidateNotifiers(); err != nil {
		return err
	}
	if err := c.ValidateFilters(); err != nil {
		return err
	}
	// Add validation for credentials
	if c.Credentials.YoutubeApiKey == "" {
		return fmt.Errorf("youtube_api_key is required")
	}
	return nil
}

// ValidateSources validates all source configurations
func (c *Config) ValidateSources() error {
	for _, src := range c.Sources {
		switch src.Type {
		case "reddit":
			if redditCfg, ok := src.Config.(*reddit.RedditCrawlerConfig); ok {
				if err := redditCfg.Validate(); err != nil {
					return fmt.Errorf("source %s: %w", src.Name, err)
				}
			}
		case "youtube":
			if youtubeCfg, ok := src.Config.(*youtube.YouTubeCrawlerConfig); ok {
				if err := youtubeCfg.Validate(); err != nil {
					return fmt.Errorf("source %s: %w", src.Name, err)
				}
			}
		default:
			return fmt.Errorf("unsupported source type: %s", src.Type)
		}
	}
	return nil
}

// ValidateNotifiers validates all notifier configurations
func (c *Config) ValidateNotifiers() error {
	for _, notifier := range c.Notifiers {
		// Attempt to validate the notifier's inner config if it supports it
		if notifier.Config == nil {
			continue
		}
		if v, ok := notifier.Config.(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return fmt.Errorf("notifier %s: %w", notifier.Name, err)
			}
		}
	}
	return nil
}

// ValidateFilters validates all filter configurations
func (c *Config) ValidateFilters() error {
	for _, filter := range c.Filters {
		if err := filter.Validate(); err != nil {
			return fmt.Errorf("filter %s: %w", filter.Name, err)
		}
	}
	return nil
}
