package feed

import "github.com/mmcdole/gofeed"

type Provider interface {
	TransformURL(rawURL string) string
	ShouldSkipItem(item *gofeed.Item) bool
}

func CreateProvider(feedType string) Provider {
	switch feedType {
	case "youtube":
		return NewYoutubeProvider()
	default:
		return NewRSSProvider()
	}
}
