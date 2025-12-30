package model

import "context"

// Source defines the interface for content sources
type Source interface {
	// Name returns the name of the source
	Name() string

	// Fetch retrieves new content since the last check
	Fetch(ctx context.Context) ([]Content, error)

	// Type returns the platform type (Reddit, YouTube, etc.)
	Type() string

	// SourceID returns the identifier of the source
	// this should be unique for each source
	SourceID() string
}
