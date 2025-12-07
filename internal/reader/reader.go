package reader

import (
	"net/http"
	"strings"
	"time"
	"html"

	"github.com/mmcdole/gofeed"
)

type Reader struct {
	Parser *gofeed.Parser
}

func NewReader(client *http.Client) Reader {
	parser := gofeed.NewParser()
	parser.Client = client
	return Reader{
		Parser: parser,
	}
}

func (reader *Reader)Parse(feedUrl string, lastCheck *time.Time, jobs chan<- string) error {
	feed, err := reader.Parser.ParseURL(feedUrl)
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
		jobs <- formatMessage(channelName, url, title)
	}
	return nil
}

func formatMessage(channelName, url, title string) string {
	var builder strings.Builder

	builder.Grow(len(channelName) + len(url) + len(title) + 50)
	builder.WriteString("FROM: <b>")
	builder.WriteString(html.EscapeString(channelName))
	builder.WriteString("</b>\n<a href=\"")
	builder.WriteString(url)
	builder.WriteString("\">")
	builder.WriteString(html.EscapeString(title))
	builder.WriteString("</a>")

	return builder.String()
}

func shouldSkip(url string, publishedAt time.Time, lastCheck *time.Time) bool {
	if strings.Contains(url, "shorts") { return true }
	if lastCheck == nil && publishedAt.Before(time.Now().Add(-24 * time.Hour)) { return true }
	if lastCheck != nil && publishedAt.Before(*lastCheck) { return true }

	return false
}
