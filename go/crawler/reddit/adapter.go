package reddit

import (
	"context"
	"time"

	"github.com/ryansiau/utilities/go/model"
)

// Adapter adapts the Reddit RSS fetcher to the Source interface
type Adapter struct {
	config *RedditCrawlerConfig
	name   string
}

// NewAdapter creates a new Reddit adapter with configuration
func NewAdapter(config *RedditCrawlerConfig, name string) (*Adapter, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	if name == "" {
		name = "Reddit: r/" + config.Subreddit
	}

	return &Adapter{
		config: config,
		name:   name,
	}, nil
}

// Name returns the name of the source
func (r *Adapter) Name() string {
	return r.name
}

// Type returns the platform type
func (r *Adapter) Type() string {
	return "reddit"
}

// Fetch retrieves new content from Reddit
func (r *Adapter) Fetch(ctx context.Context) ([]model.Content, error) {
	feed, err := FetchRSS(ctx, r.config.Subreddit)
	if err != nil {
		return nil, err
	}

	var contents []model.Content
	for _, post := range feed.Entry {
		// Convert Reddit post to generic Content
		content := model.Content{
			ID:          post.ID,
			Title:       post.Title,
			URL:         post.Link.Href,
			Author:      post.Author.Name,
			Platform:    "Reddit",
			PublishedAt: parseTime(post.Published),
			UpdatedAt:   time.Now(),
			Metadata: map[string]interface{}{
				"content": string(post.Content),
			},
		}
		contents = append(contents, content)
	}

	return contents, nil
}

// parseTime converts a string time to a time.Time
func parseTime(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		// Return current time if parsing fails
		return time.Now()
	}
	return t
}
