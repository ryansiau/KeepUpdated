package title

import (
	"fmt"
	"strings"

	"github.com/ryansiau/utilities/go/model"
)

type Config struct {
	Substring string `yaml:"substring"`
	Prefix    string `yaml:"prefix"`
	Suffix    string `yaml:"suffix"`
}

func (c *Config) Validate() error {
	if c.Substring == "" && c.Prefix == "" && c.Suffix == "" {
		return fmt.Errorf("title filter conditions are empty")
	}

	return nil
}

func (c *Config) IsFilterConfig() {}

func (c *Config) Build() (model.Filter, error) {
	return NewTitleFilter(c), nil
}

type Title struct {
	config *Config
}

var _ model.Filter = (*Title)(nil)

func NewTitleFilter(config *Config) model.Filter {
	return &Title{config: config}
}

func (t *Title) Name() string {
	return "Title"
}

func (t *Title) Apply(content model.Content) bool {
	if t.config.Substring != "" {
		if !strings.Contains(content.Title, t.config.Substring) {
			return false
		}
	}

	if t.config.Prefix != "" {
		if !strings.HasPrefix(content.Title, t.config.Prefix) {
			return false
		}
	}

	if t.config.Suffix != "" {
		if !strings.HasSuffix(content.Title, t.config.Suffix) {
			return false
		}
	}

	return true
}

func (t *Title) Type() string {
	return "Title"
}
