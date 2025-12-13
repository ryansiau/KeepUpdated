package youtube

import (
	"fmt"

	"github.com/ryansiau/utilities/go/model"
)

// Config represents the configuration for a YouTube source
type Config struct {
	ChannelID string `yaml:"channel_id" mapstructure:"channel_id"`
}

// Validate validates the YouTube source configuration
func (y *Config) Validate() error {
	if y.ChannelID == "" {
		return fmt.Errorf("channel_id is required")
	}
	return nil
}

func (y *Config) IsCrawler() {}

func (c *Config) Build(name string) (model.Source, error) {
	// TODO how to handle api key?
	return NewAdapter(c, name, "")
}
