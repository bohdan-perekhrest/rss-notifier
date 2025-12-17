package feed

import (
	"strings"
	"github.com/mmcdole/gofeed"
)

type YoutubeProvider struct {}

func NewYoutubeProvider() *YoutubeProvider {
	return &YoutubeProvider{}
}

func (y *YoutubeProvider) TransformURL(channelID string) string {
	return "https://www.youtube.com/feeds/videos.xml?channel_id=" + channelID
}

func (y *YoutubeProvider) ShouldSkipItem(item *gofeed.Item) bool {
	return strings.Contains(item.Link, "/shorts/")
}
