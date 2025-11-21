package reddit

import (
	"context"
	"encoding/xml"
	"fmt"

	"resty.dev/v3"

	"github.com/ryansiau/utilities/go/config"
)

func FetchRSS(ctx context.Context, subreddit string) (*Feed, error) {
	res := Feed{}

	client := resty.New()
	client.SetHeader("User-Agent", config.HTTPClientUserAgent)

	resp, err := client.R().
		SetContext(ctx).
		SetResult(&res).
		Get("https://www.reddit.com/r/" + subreddit + "/.rss")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	return &res, nil
}

type Feed struct {
	XMLName  xml.Name `xml:"feed"`
	Text     string   `xml:",chardata"`
	Xmlns    string   `xml:"xmlns,attr"`
	Media    string   `xml:"media,attr"`
	Category struct {
		Text  string `xml:",chardata"`
		Term  string `xml:"term,attr"`
		Label string `xml:"label,attr"`
	} `xml:"category"`
	Updated string `xml:"updated"`
	Icon    string `xml:"icon"`
	ID      string `xml:"id"`
	Link    []struct {
		Text string `xml:",chardata"`
		Rel  string `xml:"rel,attr"`
		Href string `xml:"href,attr"`
		Type string `xml:"type,attr"`
	} `xml:"link"`
	Logo     string `xml:"logo"`
	Subtitle string `xml:"subtitle"`
	Title    string `xml:"title"`
	Entry    []struct {
		Text   string `xml:",chardata"`
		Author struct {
			Text string `xml:",chardata"`
			Name string `xml:"name"`
			URI  string `xml:"uri"`
		} `xml:"author"`
		Category struct {
			Text  string `xml:",chardata"`
			Term  string `xml:"term,attr"`
			Label string `xml:"label,attr"`
		} `xml:"category"`
		Content xml.CharData `xml:"content"`
		ID      string       `xml:"id"`
		Link    struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
		} `xml:"link"`
		Updated   string `xml:"updated"`
		Published string `xml:"published"`
		Title     string `xml:"title"`
		Thumbnail struct {
			Text string `xml:",chardata"`
			URL  string `xml:"url,attr"`
		} `xml:"thumbnail"`
	} `xml:"entry"`
}
