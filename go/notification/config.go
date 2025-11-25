package notification

type NotifierConfig struct {
	Type   string      `yaml:"type"`
	Name   string      `yaml:"name"`
	Config interface{} `yaml:"config,omitempty"`
}
