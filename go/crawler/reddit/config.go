package reddit

import (
	"fmt"
)

// RedditCrawlerConfig represents the configuration for a Reddit source
type RedditCrawlerConfig struct {
	Subreddit string `yaml:"subreddit" mapstructure:"subreddit"`
}

// Validate validates the Reddit crawler configuration
func (r *RedditCrawlerConfig) Validate() error {
	if r.Subreddit == "" {
		return fmt.Errorf("subreddit is required")
	}
	return nil
}

func (r *RedditCrawlerConfig) IsCrawler() {}
