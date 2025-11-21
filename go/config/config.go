package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type CrawlerConfig interface {
	Validate() error
}

// RedditCrawlerConfig represents the configuration for a Reddit source
type RedditCrawlerConfig struct {
	Username  string `yaml:"username"`
	Subreddit string `yaml:"subreddit"`
	Interval  string `yaml:"interval"`
}

// Validate validates the Reddit crawler configuration
func (r *RedditCrawlerConfig) Validate() error {
	if r.Username == "" {
		return fmt.Errorf("username is required")
	}
	if r.Subreddit == "" {
		return fmt.Errorf("subreddit is required")
	}
	return nil
}

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
		var cfg RedditCrawlerConfig
		if raw.Config != nil {
			data, _ := yaml.Marshal(raw.Config)
			_ = yaml.Unmarshal(data, &cfg)
		}
		s.Config = &cfg
	default:
		var cfg RedditCrawlerConfig
		if raw.Config != nil {
			data, _ := yaml.Marshal(raw.Config)
			_ = yaml.Unmarshal(data, &cfg)
		}
		s.Config = &cfg
	}
	return nil
}

type NotifierConfig struct {
	Type   string      `yaml:"type"`
	Name   string      `yaml:"name"`
	Config interface{} `yaml:"config,omitempty"`
}

type FilterConfig struct {
	Name      string      `yaml:"name"`
	Type      string      `yaml:"type"`
	Config    interface{} `yaml:"config,omitempty"`
	Sources   []string    `yaml:"sources"`
	Notifiers []string    `yaml:"notifiers"`
}

type ScheduleConfig struct {
	Interval string `yaml:"interval"`
}

type Config struct {
	Sources   []SourceConfig   `yaml:"sources"`
	Notifiers []NotifierConfig `yaml:"notifiers"`
	Filters   []FilterConfig   `yaml:"filters"`
	Schedule  ScheduleConfig   `yaml:"schedule"`
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
	return nil
}

// ValidateSources validates all source configurations
func (c *Config) ValidateSources() error {
	for _, src := range c.Sources {
		if err := src.Config.Validate(); err != nil {
			return fmt.Errorf("source %s: %w", src.Name, err)
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

// Minimal Validate for FilterConfig to satisfy compile when no external validator is present
func (f *FilterConfig) Validate() error {
	if f.Name == "" {
		return fmt.Errorf("filter name is required")
	}
	if f.Type == "" {
		return fmt.Errorf("filter type is required")
	}
	return nil
}
