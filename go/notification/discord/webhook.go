package discord

import (
	"context"
	"fmt"
	"net/url"

	"resty.dev/v3"

	"github.com/ryansiau/utilities/go/model"
)

type Config struct {
	URL string
}

func (c *Config) Validate() error {
	if _, err := url.Parse(c.URL); err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	return nil
}

func (c *Config) IsNotifierConfig() {}

type DiscordWebhook struct {
	conf   *Config
	client *resty.Client
}

func NewDiscordWebhook(conf *Config) (model.Notifier, error) {
	return &DiscordWebhook{
		conf:   conf,
		client: resty.New(),
	}, nil
}

func (d *DiscordWebhook) Name() string {
	return "Discord"
}

func (d *DiscordWebhook) Type() string {
	return "Discord"
}

func (d *DiscordWebhook) Send(ctx context.Context, content model.Content) error {
	payload := WebhookPayload{
		Embeds: []Embed{
			{
				Title:       content.Title,
				Description: content.Description,
				Url:         content.URL,
				Color:       14408667,
				Author: EmbedAuthor{
					Name:    content.Author,
					Url:     "",
					IconUrl: "",
				},
				Footer: EmbedFooter{
					Text:    content.Platform,
					IconUrl: "",
				},
				Timestamp: content.PublishedAt,
				Image:     EmbedImage{},
				Fields:    nil,
			},
		},
		Attachments: nil,
	}

	req := d.client.R().SetContext(ctx)

	req.SetHeader("Content-Type", "application/json")
	req.SetBody(payload)

	resp, err := req.Post(d.conf.URL)
	if err != nil {
		return err
	}
	fmt.Println(resp.String())

	if resp.IsError() {
		return fmt.Errorf("received status code %d: %s", resp.StatusCode(), resp.String())
	}

	return nil
}
