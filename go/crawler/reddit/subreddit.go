package reddit

import (
	"context"
	"fmt"

	"resty.dev/v3"

	"github.com/ryansiau/utilities/go/config"
)

func FetchSubreddit(ctx context.Context, subreddit string) (*Subreddit, error) {
	res := Subreddit{}

	client := resty.New()
	client.SetHeader("User-Agent", config.HTTPClientUserAgent)

	resp, err := client.R().
		SetContext(ctx).
		SetResult(&res).
		Get("https://www.reddit.com/r/" + subreddit + ".json")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	return &res, nil
}

type Subreddit struct {
	Kind string `json:"kind"`
	Data struct {
		Children []struct {
			Kind string `json:"kind"`
			Data Post   `json:"data"`
		} `json:"children"`
		After   string `json:"after"`
		Before  string `json:"before"`
		Dist    int    `json:"dist"`
		Modhash string `json:"modhash"`
	} `json:"data"`
}

type Post struct {
	Title       string  `json:"title"`
	Subreddit   string  `json:"subreddit"`
	Author      string  `json:"author"`
	URL         string  `json:"url"`
	Permalink   string  `json:"permalink"`
	Score       int     `json:"score"`
	NumComments int     `json:"num_comments"`
	CreatedUTC  float64 `json:"created_utc"`
	SelfText    string  `json:"selftext"`
	Thumbnail   string  `json:"thumbnail"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	UpvoteRatio float64 `json:"upvote_ratio"`
	Stickied    bool    `json:"stickied"`
	Over18      bool    `json:"over_18"`
	IsSelf      bool    `json:"is_self"`
}
