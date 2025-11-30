package metadata

import (
	"fmt"
	"strings"

	"github.com/ryansiau/utilities/go/model"
)

type Config struct {
	Conditions []FilterCondition `yaml:"conditions"`
}

type FilterCondition struct {
	Comp  string      `yaml:"comp"`
	Field string      `yaml:"field"`
	Value interface{} `yaml:"value"`
}

// Comparators
const (
	// for comparables
	CompEqual    = "equal"
	CompNotEqual = "not_equal"

	// for string
	CompContains    = "contains"
	CompNotContains = "not_contains"
)

func validateComp(comp string) bool {
	return comp == CompEqual ||
		comp == CompNotEqual ||
		comp == CompContains ||
		comp == CompNotContains
}

func (c *Config) Validate() error {
	for _, filter := range c.Conditions {
		if validateComp(filter.Comp) {
			return fmt.Errorf("comp value is unknown or unsupported")
		}
		if filter.Field == "" {
			return fmt.Errorf("field name is empty")
		}
		if filter.Value == "" {
			return fmt.Errorf("value is empty")
		}

		if filter.Comp == CompContains || filter.Comp == CompNotContains {
			_, ok := filter.Value.(string)
			if !ok {
				return fmt.Errorf("value has to be a string for the selected comp")
			}
		}
	}

	return nil
}

func (c *Config) IsFilterConfig() {}

type Metadata struct {
	config *Config
}

func NewMetadata(config *Config) model.Filter {
	return &Metadata{config}
}

func (m Metadata) Name() string {
	return "Metadata"
}

func (m Metadata) Apply(content model.Content) bool {
	for _, cond := range m.config.Conditions {
		switch cond.Comp {
		case CompEqual:
			if content.Metadata[cond.Field] != cond.Value {
				return false
			}
		case CompNotEqual:
			if content.Metadata[cond.Field] == cond.Value {
				return false
			}
		case CompContains:
			val, ok := content.Metadata[cond.Field].(string)
			if !ok {
				return false
			}
			if !strings.Contains(val, cond.Value.(string)) {
				return false
			}
		case CompNotContains:
			val, ok := content.Metadata[cond.Field].(string)
			if !ok {
				return false
			}
			if strings.Contains(val, cond.Value.(string)) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func (m Metadata) Type() string {
	return "Metadata"
}
