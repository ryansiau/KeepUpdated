package generic_rss

import (
	"fmt"
	"net/url"

	"github.com/ryansiau/utilities/go/model"
)

// Config holds RSS source configuration
type Config struct {
	FeedURL string `json:"feed_url"`
}

func (c *Config) Validate() error {
	if c.FeedURL == "" {
		return fmt.Errorf("feed_url is required")
	}

	_, err := url.Parse(c.FeedURL)
	if err != nil {
		return fmt.Errorf("failed to parse feed URL: %w", err)
	}
	return nil
}

func (c *Config) IsCrawler() {}

func (c *Config) Build(name string) (model.Source, error) {
	return New(c, name), nil
}
