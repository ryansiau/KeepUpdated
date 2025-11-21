package reddit

import (
	"context"
	"time"

	"github.com/ryansiau/utilities/go/crawler"
)

// Adapter adapts the Reddit RSS fetcher to the Source interface
type Adapter struct {
	subreddit string
	name      string
}

// NewAdapter creates a new Reddit adapter
func NewAdapter(subreddit string, name string) *Adapter {
	if name == "" {
		name = "Reddit: r/" + subreddit
	}
	return &Adapter{
		subreddit: subreddit,
		name:      name,
	}
}

// Name returns the name of the source
func (r *Adapter) Name() string {
	return r.name
}

// Type returns the platform type
func (r *Adapter) Type() string {
	return "Reddit"
}

// Fetch retrieves new content from Reddit
func (r *Adapter) Fetch(ctx context.Context) ([]crawler.Content, error) {
	feed, err := FetchRSS(ctx, r.subreddit)
	if err != nil {
		return nil, err
	}

	var contents []crawler.Content
	for _, post := range feed.Entry {
		// Convert Reddit post to generic Content
		content := crawler.Content{
			ID:          post.ID,
			Title:       post.Title,
			URL:         post.Link.Href,
			Author:      post.Author.Name,
			Platform:    "Reddit",
			PublishedAt: parseTime(post.Published),
			UpdatedAt:   parseTime(post.Updated),
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
