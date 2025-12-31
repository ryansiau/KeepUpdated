package filter

import (
	"strings"

	"github.com/ryansiau/KeepUpdated/go/model"
)

// KeywordFilter is a filter that checks for keywords in the content title
type KeywordFilter struct {
	name     string
	keywords []string
}

// NewKeywordFilter creates a new keyword filter
func NewKeywordFilter(name string, keywords []string) *KeywordFilter {
	if name == "" {
		name = "Keyword Filter"
	}
	return &KeywordFilter{
		name:     name,
		keywords: keywords,
	}
}

// Name returns the name of the filter
func (f *KeywordFilter) Name() string {
	return f.name
}

// Type returns the filter type
func (f *KeywordFilter) Type() string {
	return "Keyword"
}

// Apply checks if the content passes the filter
func (f *KeywordFilter) Apply(content model.Content) bool {
	// If no keywords are specified, all content passes the filter
	if len(f.keywords) == 0 {
		return true
	}

	// Check if any keyword is in the title
	for _, keyword := range f.keywords {
		if containsIgnoreCase(content.Title, keyword) {
			return true
		}
	}

	return false
}

// containsIgnoreCase checks if a string contains a substring, ignoring case
func containsIgnoreCase(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}
