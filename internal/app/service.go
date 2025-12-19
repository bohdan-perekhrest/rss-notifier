package app

import (
	"fmt"
	"log"
	"sync"

	"rss-notifier/config"
	"rss-notifier/internal/cache"
	"rss-notifier/internal/feed"
	"rss-notifier/internal/notifier"
	"rss-notifier/internal/reader"
)

type Service struct {
	reader   *reader.Reader
	notifier notifier.Notifier
	cache    cache.Cache
	cfg      *config.Config
}

func NewService(
	reader *reader.Reader,
	notifier notifier.Notifier,
	cache cache.Cache,
	cfg *config.Config,
) *Service {
	return &Service{
		reader:   reader,
		notifier: notifier,
		cache:    cache,
		cfg:      cfg,
	}
}

func (s *Service) Run() error {
	feedConfigs, err := s.cfg.LoadFeeds()
	if err != nil {
		return fmt.Errorf("failed to load feeds: %w", err)
	}

	feeds := make([]feed.Feed, len(feedConfigs))
	for i, fc := range feedConfigs {
		feeds[i] = feed.Feed{
			URL:      fc.URL,
			Type:     fc.Type,
			Provider: feed.CreateProvider(fc.Type),
		}
	}

	responses := make(chan *feed.Response, len(feeds))
	var wg sync.WaitGroup

	for _, f := range feeds {
		wg.Add(1)
		go func(f feed.Feed) {
			defer wg.Done()

			resp, err := s.reader.Parse(f)
			if err != nil {
				log.Printf("Failed to parse feed %s: %v", f.URL, err)
				return
			}

			if len(resp.Items) > 0 {
				responses <- resp
			}
		}(f)
	}

	go func() {
		wg.Wait()
		close(responses)
	}()

	var results []feed.Response
	for resp := range responses {
		results = append(results, *resp)
	}

	if len(results) == 0 {
		log.Println("No new items found")
		return nil
	}

	message := notifier.FormatFeeds(results)
	if err := s.notifier.Send(message); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	cacheData := make(map[string]string, len(results))
	for _, resp := range results {
		if resp.NewestItemURL != "" {
			cacheData[resp.FeedURL] = resp.NewestItemURL
		}
	}

	if err := s.cache.SaveLastItemURL(cacheData); err != nil {
		return fmt.Errorf("failed to update cache: %w", err)
	}

	log.Printf("Successfully processed %d feeds with new items", len(results))
	return nil
}
