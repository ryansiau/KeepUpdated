package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the top-level configuration structure for the system.
type Config struct {
	Sources   []SourceConfig   `yaml:"sources"`
	Notifiers []NotifierConfig `yaml:"notifiers"`
	Filters   []FilterConfig   `yaml:"filters"`
	Schedule  ScheduleConfig   `yaml:"schedule"`
}

// SourceConfig represents the configuration for a content source.
type SourceConfig struct {
	Type   string      `yaml:"type"`
	Name   string      `yaml:"name"`
	Config interface{} `yaml:"config,omitempty"`
}

// NotifierConfig represents the configuration for a notification channel.
type NotifierConfig struct {
	Type   string      `yaml:"type"`
	Name   string      `yaml:"name"`
	Config interface{} `yaml:"config,omitempty"`
}

// FilterConfig represents the configuration for a content filter.
type FilterConfig struct {
	Name      string      `yaml:"name"`
	Type      string      `yaml:"type"`
	Config    interface{} `yaml:"config,omitempty"`
	Sources   []string    `yaml:"sources"`
	Notifiers []string    `yaml:"notifiers"`
}

// ScheduleConfig represents the configuration for scheduling.
type ScheduleConfig struct {
	Interval string `yaml:"interval"`
}

func LoadConfig(filepath string) (*Config, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	return &cfg, nil
}

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
