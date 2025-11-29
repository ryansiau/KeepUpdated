package model

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
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
	Metadata    Metadata `gorm:"type:text"`
}

type Metadata map[string]interface{}

func (m Metadata) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	data, err := json.Marshal(m)
	return string(data), err
}

func (m *Metadata) Scan(value interface{}) error {
	if value == nil {
		*m = make(map[string]interface{})
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("invalid type for MetadataMap")
	}
	err := json.Unmarshal(bytes, m)
	return err
}

func (m Metadata) GormDataType() string {
	return "text"
}
