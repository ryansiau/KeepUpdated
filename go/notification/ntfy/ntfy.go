package ntfy

import (
	"context"
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"
	"resty.dev/v3"

	"github.com/ryansiau/utilities/go/common"
	"github.com/ryansiau/utilities/go/model"
)

type Config struct {
	BaseURL     string `yaml:"base_url"`
	Priority    string `yaml:"priority"`
	AccessToken string `yaml:"access_token"`
	Topic       string `yaml:"topic"`
}

func (c *Config) Validate() error {
	if _, err := url.Parse(c.BaseURL); err != nil {
		return err
	}
	if c.Topic == "" {
		return fmt.Errorf("topic is required")
	}
	if c.AccessToken == "" {
		logrus.Warn("no access token provided, notifications to ntfy might fail if the host requires authentication")
	}

	return nil
}

func (c *Config) IsNotifierConfig() {}

type Ntfy struct {
	conf   *Config
	client *resty.Client
}

func New(conf *Config) (model.Notifier, error) {
	return &Ntfy{
		conf:   conf,
		client: resty.New().SetHeader("User-Agent", common.HTTPClientUserAgent),
	}, nil
}

func (n *Ntfy) Name() string {
	return "ntfy"
}

func (n *Ntfy) Send(ctx context.Context, content model.Content) error {
	req := n.client.R().SetContext(ctx)

	req.SetHeader("Accept", "application/json")
	req.SetAuthScheme("Bearer")
	req.SetAuthToken(n.conf.AccessToken)

	req.SetHeader("X-Message", content.Title)
	req.SetHeader("X-Title", fmt.Sprintf("New update from %s (%s)", content.Author, content.Platform))
	req.SetHeader("X-Priority", n.conf.Priority)
	req.SetHeader("X-Click", content.URL)

	finalURL, err := url.JoinPath(n.conf.BaseURL, n.conf.Topic)
	if err != nil {
		return err
	}

	resp, err := req.Post(finalURL)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("received status code %d: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

func (n *Ntfy) Type() string {
	return "ntfy"
}
