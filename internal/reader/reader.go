package reader

import (
	"github.com/mmcdole/gofeed"
	"log"
	"rss-notifier/internal/telegram"
	"strings"
	"time"
)

type Reader struct {}

func NewReader() *Reader {
	return &Reader{}
}

func (reader *Reader)Parse(feedUrl string, lastCheck *time.Time) error {
	feed, err := fetchFeed(feedUrl)
	if err != nil {
		return err
	}

	channelName := feed.Title
	for _, item := range feed.Items {
		url := item.Link
		publishedAt, err := time.Parse(time.RFC3339, item.Published)
		if err != nil {
			continue
		}

		if shouldSkip(url, publishedAt, lastCheck) {
			continue
		}
		title := item.Title
		description := item.Description
		message := formatMessage(channelName, url, title, description)
		if err := telegram.SendMessage(message); err != nil { log.Printf("Failed to send message: %v\n", err) }
	}
	return nil
}

func fetchFeed(url string) (*gofeed.Feed, error) {
	parser := gofeed.NewParser()
	return parser.ParseURL(url)
}

func formatMessage(channelName, url, title, description string) string {
	return "FROM: <b>" + escapeHTML(channelName) + "</b>\nDescription:" + escapeHTML(description) + "\n<a href=\"" + url + "\">" + escapeHTML(title) + "</a>"
}

func escapeHTML(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	return text
}

func shouldSkip(url string, publishedAt time.Time, lastCheck *time.Time) bool {
	if strings.Contains(url, "shorts") { return true }
	if lastCheck == nil && publishedAt.Before(time.Now().Add(-24 * time.Hour)) { return true }
	if lastCheck != nil && publishedAt.Before(*lastCheck) { return true }

	return false
}
