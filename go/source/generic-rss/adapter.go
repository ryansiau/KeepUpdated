package generic_rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"time"

	"resty.dev/v3"

	"github.com/ryansiau/KeepUpdated/go/common"
	"github.com/ryansiau/KeepUpdated/go/model"
)

// RSSFeed represents the root RSS structure
type RSSFeed struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Content     string `xml:"encoded"` // Often used for full content
	PubDate     string `xml:"pubDate"` // Publication date as string
	GUID        string `xml:"guid"`    // Unique identifier
	Author      string `xml:"author"`  // Optional author field
	Creator     string `xml:"creator"` // Dublin Core creator
}

// Adapter implements the Source interface for RSS feeds
type Adapter struct {
	feedURL string
	name    string
	client  *resty.Client
}

// New creates a new RSS source
func New(config *Config, name string) model.Source {
	client := resty.New().
		//SetTimeout(5*time.Second).
		SetHeader("User-Agent", common.HTTPClientUserAgent)
	return &Adapter{
		feedURL: config.FeedURL,
		name:    name,
		client:  client,
	}
}

// Name returns the name of the source
func (r *Adapter) Name() string {
	return r.name
}

// Type returns the platform type
func (r *Adapter) Type() string {
	return "RSS"
}

// Fetch retrieves new content since the last check
func (r *Adapter) Fetch(ctx context.Context) ([]model.Content, error) {
	resp, err := r.client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/rss+xml").
		SetHeader("Accept", "application/xml").
		SetHeader("Accept", "text/xml").
		SetHeader("User-Agent", common.HTTPClientUserAgent).
		Get(r.feedURL)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("RSS feed returned status: %d %s", resp.StatusCode(), resp.Status())
	}

	content, err := r.parseFeed(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	return content, nil
}

// parseFeed parses the RSS feed and returns new content
func (r *Adapter) parseFeed(reader io.Reader) ([]model.Content, error) {
	var feed RSSFeed
	if err := xml.NewDecoder(reader).Decode(&feed); err != nil {
		return nil, fmt.Errorf("failed to decode XML: %w", err)
	}

	var contents []model.Content

	for _, item := range feed.Channel.Items {
		// Use GUID as unique identifier, fall back to link
		itemID := item.GUID
		if itemID == "" {
			itemID = item.Link
		}

		// Parse publication date
		pubDate, err := parseDate(item.PubDate)
		if err != nil {
			// If we can't parse the date, use the current time
			pubDate = time.Now()
		}

		// Determine author
		author := item.Author
		if author == "" {
			author = item.Creator
		}
		if author == "" {
			author = r.name // Fall back to source name
		}

		// Determine content text
		contentText := item.Description
		if item.Content != "" {
			contentText = item.Content
		}

		content := model.Content{
			ID:          strings.TrimSpace(itemID),
			SourceID:    r.SourceID(),
			Title:       strings.TrimSpace(item.Title),
			Description: strings.TrimSpace(contentText),
			URL:         strings.TrimSpace(item.Link),
			Author:      strings.TrimSpace(author),
			Platform:    strings.TrimSpace(r.name),
			PublishedAt: pubDate,
			UpdatedAt:   time.Now(),
			Metadata:    nil,
		}

		contents = append(contents, content)
	}

	return contents, nil
}

func (r *Adapter) SourceID() string {
	return fmt.Sprintf("RSS:%s", r.feedURL)
}

// parseDate attempts to parse various RSS date formats
func parseDate(dateStr string) (time.Time, error) {
	// Try common RSS date formats
	formats := []string{
		time.RFC1123,  // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC1123Z, // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC822,   // "02 Jan 06 15:04 MST"
		time.RFC822Z,  // "02 Jan 06 15:04 -0700"
		time.RFC3339,  // "2006-01-02T15:04:05Z07:00"
		"02 Jan 2006 15:04:05 MST",
		"Mon, 2 Jan 2006 15:04:05 -0700",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
