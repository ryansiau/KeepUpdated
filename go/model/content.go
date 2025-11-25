package model

import (
	"context"
	"time"
)

// Source defines the interface for content sources
type Source interface {
	// Name returns the name of the source
	Name() string

	// Fetch retrieves new content since the last check
	Fetch(ctx context.Context) ([]Content, error)

	// Type returns the platform type (Reddit, YouTube, etc.)
	Type() string
}

// Content represents a generic content item from any platform
type Content struct {
	ID          string
	Title       string
	Description string
	URL         string
	Author      string
	Platform    string
	PublishedAt time.Time
	UpdatedAt   time.Time
	Metadata    map[string]interface{}
}
