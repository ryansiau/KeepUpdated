// Package twitter is supposed to be a package that handles crawling user's new posts from twitter Official API.
// However, twitter's public API is very limited and does not seem to allow even 4 API calls/day.
// This means the API key has to be of a paid plan to allow more calls, and I'm not fond of the idea
// to get paid plan just to get updates from some handles.
// Therefore, this package is held and might be canceled.
package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Content represents a tweet/post from Twitter
type Content struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Author      string    `json:"author"`
	Text        string    `json:"text"`
	PublishedAt time.Time `json:"published_at"`
	MediaURLs   []string  `json:"media_urls,omitempty"`
	Likes       int       `json:"likes,omitempty"`
	Retweets    int       `json:"retweets,omitempty"`
	Replies     int       `json:"replies,omitempty"`
	QuoteTweets int       `json:"quote_tweets,omitempty"`
}

// Config holds Twitter API configuration for a specific user
type Config struct {
	BearerToken     string `json:"bearer_token"`
	Username        string `json:"username"`         // The Twitter username to monitor
	UserID          string `json:"user_id"`          // Optional: User ID (more reliable than username)
	SinceID         string `json:"since_id"`         // Last fetched tweet ID
	MaxResults      int    `json:"max_results"`      // Max tweets per request (1-100)
	IncludeRetweets bool   `json:"include_retweets"` // Whether to include retweets
	IncludeReplies  bool   `json:"include_replies"`  // Whether to include replies
}

// Twitter implements the Source interface for Twitter/X user posts
type Twitter struct {
	config Config
	client *http.Client
	name   string
}

// New creates a new Twitter crawler for a specific user
func New(config Config) *Twitter {
	name := fmt.Sprintf("Twitter Crawler - @%s", config.Username)
	if config.Username == "" && config.UserID != "" {
		name = fmt.Sprintf("Twitter Crawler - User %s", config.UserID)
	}

	return &Twitter{
		config: config,
		client: &http.Client{Timeout: 30 * time.Second},
		name:   name,
	}
}

// Name returns the name of the source
func (t *Twitter) Name() string {
	return t.name
}

// Type returns the platform type
func (t *Twitter) Type() string {
	return "Twitter"
}

// Fetch retrieves new posts from the user since the last check
func (t *Twitter) Fetch(ctx context.Context) ([]Content, error) {
	userID := t.config.UserID

	// If we don't have UserID, resolve username to UserID first
	if userID == "" && t.config.Username != "" {
		resolvedUserID, err := t.resolveUsername(ctx, t.config.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve username %s: %w", t.config.Username, err)
		}
		userID = resolvedUserID
		t.config.UserID = userID // Cache for future use
	}

	if userID == "" {
		return nil, fmt.Errorf("no user ID or username provided")
	}

	return t.fetchUserTweets(ctx, userID)
}

// resolveUsername converts a username to a User ID
func (t *Twitter) resolveUsername(ctx context.Context, username string) (string, error) {
	url := fmt.Sprintf("https://api.twitter.com/2/users/by/username/%s", username)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.config.BearerToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("twitter API returned status: %d for username resolution", resp.StatusCode)
	}

	var result struct {
		Data struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data.ID, nil
}

// fetchUserTweets fetches tweets from a specific user by their ID
func (t *Twitter) fetchUserTweets(ctx context.Context, userID string) ([]Content, error) {
	// Build query parameters
	params := fmt.Sprintf("max_results=%d", t.config.MaxResults)

	if t.config.SinceID != "" {
		params += fmt.Sprintf("&since_id=%s", t.config.SinceID)
	}

	// Build exclusions
	var exclusions []string
	if !t.config.IncludeRetweets {
		exclusions = append(exclusions, "retweets")
	}
	if !t.config.IncludeReplies {
		exclusions = append(exclusions, "replies")
	}

	if len(exclusions) > 0 {
		params += "&exclude="
		for i, exclusion := range exclusions {
			if i > 0 {
				params += ","
			}
			params += exclusion
		}
	}

	// Add tweet fields
	params += "&tweet.fields=author_id,created_at,public_metrics,text,conversation_id"

	url := fmt.Sprintf("https://api.twitter.com/2/users/%s/tweets?%s", userID, params)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.config.BearerToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitter API returned status: %d", resp.StatusCode)
	}

	var result twitterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return t.parseTweets(result.Data), nil
}

// parseTweets converts Twitter API response to Content objects
func (t *Twitter) parseTweets(tweets []tweetData) []Content {
	var content []Content

	for _, tweet := range tweets {
		publishedAt, _ := time.Parse(time.RFC3339, tweet.CreatedAt)

		// Determine if this is a reply
		isReply := tweet.ConversationID != "" && tweet.ConversationID != tweet.AuthorID

		c := Content{
			ID:          tweet.ID,
			Title:       fmt.Sprintf("Tweet by @%s", t.config.Username),
			URL:         fmt.Sprintf("https://twitter.com/%s/status/%s", t.config.Username, tweet.ID),
			Author:      t.config.Username,
			Text:        tweet.Text,
			PublishedAt: publishedAt,
			Likes:       tweet.PublicMetrics.LikeCount,
			Retweets:    tweet.PublicMetrics.RetweetCount,
			Replies:     tweet.PublicMetrics.ReplyCount,
			QuoteTweets: tweet.PublicMetrics.QuoteCount,
		}

		// Add context to title for replies
		if isReply {
			c.Title = fmt.Sprintf("Reply by @%s", t.config.Username)
		}

		content = append(content, c)
	}

	return content
}

// Twitter API response structures
type twitterResponse struct {
	Data []tweetData `json:"data"`
	Meta struct {
		OldestID    string `json:"oldest_id"`
		NewestID    string `json:"newest_id"`
		ResultCount int    `json:"result_count"`
	} `json:"meta"`
}

type tweetData struct {
	ID             string        `json:"id"`
	Text           string        `json:"text"`
	AuthorID       string        `json:"author_id"`
	ConversationID string        `json:"conversation_id"`
	CreatedAt      string        `json:"created_at"`
	PublicMetrics  publicMetrics `json:"public_metrics"`
}

type publicMetrics struct {
	LikeCount    int `json:"like_count"`
	RetweetCount int `json:"retweet_count"`
	ReplyCount   int `json:"reply_count"`
	QuoteCount   int `json:"quote_count"`
}

// SetSinceID allows updating the since_id parameter
func (t *Twitter) SetSinceID(sinceID string) {
	t.config.SinceID = sinceID
}

// GetSinceID returns the current since_id
func (t *Twitter) GetSinceID() string {
	return t.config.SinceID
}

// GetUserID returns the resolved User ID
func (t *Twitter) GetUserID() string {
	return t.config.UserID
}
