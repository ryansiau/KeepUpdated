package reddit

import (
	"context"
	"fmt"
	"time"

	"resty.dev/v3"

	"github.com/ryansiau/utilities/go/common"
	"github.com/ryansiau/utilities/go/model"
)

// Adapter adapts the Reddit RSS fetcher to the Source interface
type Adapter struct {
	client *resty.Client
	config *Config
	name   string
}

// NewAdapter creates a new Reddit adapter with configuration
func NewAdapter(config *Config, name string) (*Adapter, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	if name == "" {
		name = "Reddit: r/" + config.Subreddit
	}

	client := resty.New().
		SetTimeout(time.Second*5).
		SetHeader("User-Agent", common.HTTPClientUserAgent)

	return &Adapter{
		config: config,
		name:   name,
		client: client,
	}, nil
}

// Name returns the name of the source
func (a *Adapter) Name() string {
	return a.name
}

// Type returns the platform type
func (a *Adapter) Type() string {
	return "reddit"
}

// Fetch retrieves new content from Reddit
func (a *Adapter) Fetch(ctx context.Context) ([]model.Content, error) {
	feed, err := a.FetchRSS(ctx, a.config.Subreddit)
	if err != nil {
		return nil, err
	}

	var contents []model.Content
	for _, post := range feed.Entry {
		// Convert Reddit post to generic Content
		content := model.Content{
			ID:          post.ID,
			SourceID:    a.SourceID(),
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

func (a *Adapter) SourceID() string {
	return fmt.Sprintf("Reddit:%s", a.config.Subreddit)
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
