package discord

import "time"

type WebhookPayload struct {
	Content     string        `json:"content"`
	Embeds      []Embed       `json:"embeds"`
	Attachments []interface{} `json:"attachments"`
}

type Embed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Url         string       `json:"url,omitempty"`
	Color       int          `json:"color"`
	Author      EmbedAuthor  `json:"author,omitempty"`
	Footer      EmbedFooter  `json:"footer,omitempty"`
	Timestamp   time.Time    `json:"timestamp,omitempty"`
	Image       EmbedImage   `json:"image,omitempty"`
	Fields      []EmbedField `json:"fields,omitempty"`
}

type EmbedAuthor struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	IconUrl string `json:"icon_url"`
}

type EmbedFooter struct {
	Text    string `json:"text"`
	IconUrl string `json:"icon_url"`
}

type EmbedImage struct {
	Url string `json:"url"`
}

type EmbedField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
