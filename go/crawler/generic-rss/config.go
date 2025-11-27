package generic_rss

// Config holds RSS source configuration
type Config struct {
	FeedURL string `json:"feed_url"`
	Name    string `json:"name"`
}
