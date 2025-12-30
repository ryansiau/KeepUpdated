package config

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/ryansiau/utilities/go/filter"
	"github.com/ryansiau/utilities/go/notification"
	"github.com/ryansiau/utilities/go/pkg/database"
	"github.com/ryansiau/utilities/go/source"
)

// Config represents the entire configuration
type Config struct {
	Defaults  DefaultConfigs  `yaml:"defaults"`
	Workflows []Workflow      `yaml:"workflows"`
	Database  database.Config `yaml:"database"`
}

type DefaultConfigs struct {
	Interval    time.Duration             `yaml:"interval"`
	Credentials DefaultCreds              `yaml:"credentials"`
	Notifiers   []notification.BaseConfig `yaml:"notifiers"`
}

type DefaultCreds struct {
	YoutubeAPIKey string `yaml:"youtube_api_key"`
}

type Workflow struct {
	Name      string                    `yaml:"name"`
	Enabled   bool                      `yaml:"enabled"`
	Interval  time.Duration             `yaml:"interval"`
	Source    source.BaseConfig         `yaml:"source"`
	Filters   []filter.BaseConfig       `yaml:"filters"`
	Notifiers []notification.BaseConfig `yaml:"notifiers"`
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

	cfg.applyDefaults()

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
	if c.Defaults.Interval == 0 {
		logrus.Warn("Default interval is not set. Automatically setting it to daily")
		c.Defaults.Interval = 24 * time.Hour
	}

	for _, n := range c.Defaults.Notifiers {
		if err := n.Validate(); err != nil {
			return fmt.Errorf("invalid notifier config: %w", err)
		}
	}

	for _, w := range c.Workflows {
		if err := w.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (w *Workflow) Validate() error {
	if err := w.ValidateSources(); err != nil {
		return err
	}
	if err := w.ValidateNotifiers(); err != nil {
		return err
	}
	if err := w.ValidateFilters(); err != nil {
		return err
	}
	return nil
}

// ValidateSources validates all source configurations
func (w *Workflow) ValidateSources() error {
	if err := w.Source.Validate(); err != nil {
		return err
	}

	return nil
}

// ValidateFilters validates all filter configurations
func (w *Workflow) ValidateFilters() error {
	// workflows are allowed to have no filter

	for _, filter := range w.Filters {
		if err := filter.Validate(); err != nil {
			return fmt.Errorf("filter %s: %w", filter.Name, err)
		}
	}

	return nil
}

// ValidateNotifiers validates all notifier configurations
func (w *Workflow) ValidateNotifiers() error {
	// a workflow MUST HAVE at least a notifier.
	if len(w.Notifiers) == 0 {
		return fmt.Errorf("notifier is not set")
	}

	for _, notifier := range w.Notifiers {
		if err := notifier.Validate(); err != nil {
			return fmt.Errorf("notifier %s: %w", notifier.Name, err)
		}
	}

	return nil
}

func (c *Config) applyDefaults() {
	for widx, w := range c.Workflows {
		if w.Interval == 0 {
			w.Interval = c.Defaults.Interval
		}

		// replace workflows[].notifiers where type == "default" with defaults.notifiers
		var notifiers []notification.BaseConfig
		for _, n := range w.Notifiers {
			if n.Type == "default" {
				for _, d := range c.Defaults.Notifiers {
					notifiers = append(notifiers, d)
				}
				break
			} else {
				notifiers = append(notifiers, n)
			}
		}

		// if workflows[].notifiers is still empty, attempt to append with default.notifiers
		if len(notifiers) == 0 {
			notifiers = c.Defaults.Notifiers
		}

		c.Workflows[widx].Notifiers = notifiers
	}
}
