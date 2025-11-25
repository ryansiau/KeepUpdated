package youtube

import (
	"context"
	"fmt"
	"time"

	"github.com/ryansiau/utilities/go/model"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Adapter struct {
	name      string
	channelID string
	client    *youtube.Service
}

func NewAdapter(ctx context.Context, name, apiKey string, conf *YouTubeCrawlerConfig) (model.Source, error) {
	// Log API key configuration
	if apiKey == "" {
		return nil, fmt.Errorf("youtube_api_key is required")
	}

	client, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create YouTube client: %w", err)
	}

	if name == "" {
		name = "Youtube: " + conf.ChannelID
	}

	return &Adapter{
		name:      name,
		channelID: conf.ChannelID,
		client:    client,
	}, nil
}

func (a *Adapter) FetchVideos(channelID string) ([]*youtube.Video, error) {
	call := a.client.Search.List([]string{"snippet"})
	call.ChannelId(channelID)
	call.MaxResults(50)

	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch videos: %w", err)
	}

	var videos []*youtube.Video
	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			videoID := item.Id.VideoId
			if videoID != "" {
				videos = append(videos, &youtube.Video{
					Id: videoID,
					Snippet: &youtube.VideoSnippet{
						ChannelId:    item.Snippet.ChannelId,
						ChannelTitle: item.Snippet.ChannelTitle,
						Description:  item.Snippet.Description,
						PublishedAt:  item.Snippet.PublishedAt,
						Title:        item.Snippet.Title,
					},
				})
			}
		}
	}

	return videos, nil
}

func (a *Adapter) Name() string {
	return a.name
}

func (a *Adapter) Type() string {
	return "youtube"
}

func (a *Adapter) Fetch(ctx context.Context) ([]model.Content, error) {
	videos, err := a.FetchVideos(a.channelID)
	if err != nil {
		return nil, err
	}

	var contents []model.Content
	for _, video := range videos {
		publishedAt, err := time.Parse(video.Snippet.PublishedAt, time.RFC3339)
		if err != nil {
			publishedAt = time.Now()
		}
		contents = append(contents, model.Content{
			ID:          video.Id,
			Title:       video.Snippet.Title,
			Description: video.Snippet.Description,
			URL:         fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.Id),
			Author:      video.Snippet.ChannelTitle,
			Platform:    "YouTube",
			PublishedAt: publishedAt,
			UpdatedAt:   time.Now(),
			Metadata: map[string]interface{}{
				"video_id": video.Id,
			},
		})
	}

	return contents, nil
}
