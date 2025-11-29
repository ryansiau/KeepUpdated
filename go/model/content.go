package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

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
