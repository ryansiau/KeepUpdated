package youtube

import (
	"fmt"
)

// YouTubeCrawlerConfig represents the configuration for a YouTube source
type YouTubeCrawlerConfig struct {
	ChannelID string `yaml:"channel_id" mapstructure:"channel_id"`
}

// Validate validates the YouTube crawler configuration
func (y *YouTubeCrawlerConfig) Validate() error {
	if y.ChannelID == "" {
		return fmt.Errorf("channel_id is required")
	}
	return nil
}

func (y *YouTubeCrawlerConfig) IsCrawler() {}
