package filter

import (
	"github.com/ryansiau/utilities/go/crawler"
)

// Filter defines the interface for content filters
type Filter interface {
	// Name returns the name of the filter
	Name() string
	
	// Apply checks if the content passes the filter
	Apply(content crawler.Content) bool
	
	// Type returns the filter type (Keyword, Regex, etc.)
	Type() string
}