package reddit

import (
	"fmt"
)

// Config represents the configuration for a Reddit source
type Config struct {
	Subreddit string `yaml:"subreddit" mapstructure:"subreddit"`
}

// Validate validates the Reddit source configuration
func (r *Config) Validate() error {
	if r.Subreddit == "" {
		return fmt.Errorf("subreddit is required")
	}
	return nil
}

func (r *Config) IsCrawler() {}

func (c *Config) Build(name string) (model.Source, error) {
	return New(c, name), nil
}
