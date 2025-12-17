package feed

import "github.com/mmcdole/gofeed"

type RSSProvider struct {}

func NewRSSProvider() *RSSProvider {
	return &RSSProvider{}
}

func (r *RSSProvider) TransformURL(rawURL string) string {
	return rawURL
}

func (r *RSSProvider) ShouldSkipItem(item *gofeed.Item) bool {
	return false
}
