package reader

import(
	"fmt"
	"net/http"
	"time"

	"rss-notifier/internal/cache"
	"rss-notifier/internal/feed"

	"github.com/mmcdole/gofeed"
)

type Reader struct {
	parser *gofeed.Parser
	cache cache.Cache
	filter *feed.Filter
}

func New(client *http.Client, cache cache.Cache, maxAge time.Duration) *Reader {
	parser := gofeed.NewParser()
	parser.Client = client

	return &Reader{
		parser: parser,
		cache: cache,
		filter: feed.NewFilter(maxAge),
	}
}

func (r *Reader) Parse(f feed.Feed) (*feed.Response, error) {
	feedURL := f.Provider.TransformURL(f.URL)

	parsedFeed, err := r.parser.ParseURL(feedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed %s: %v", feedURL, err)
	}

	lastURL, err := r.cache.GetLastItemURL(feedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache for %s: %v", feedURL, err)
	}

	var items []feed.Item
	var newestURL string

	for _, item := range parsedFeed.Items {
		if newestURL == "" {
			newestURL = item.Link
		}

		if item.Link == lastURL {
			break
		}

		if f.Provider.ShouldSkipItem(item) {
			continue
		}

		if !r.filter.ShouldInclude(item) {
			continue
		}

		items = append(items, feed.Item{
			URL:       item.Link,
			Title:     item.Title,
			Published: *item.PublishedParsed,
		})
	}

	return &feed.Response{
		FeedURL:       feedURL,
		ChannelName:   parsedFeed.Title,
		Items:         items,
		NewestItemURL: newestURL,
	}, nil
}
