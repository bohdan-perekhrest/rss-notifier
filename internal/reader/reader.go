package reader

import (
	"html"
	"net/http"
	"strings"

	"rss-notifier/internal/database"
	"github.com/mmcdole/gofeed"
)

type Reader struct {
	Parser *gofeed.Parser
	db *database.Database
}

func NewReader(client *http.Client, database *database.Database) Reader {
	parser := gofeed.NewParser()
	parser.Client = client
	return Reader{
		Parser: parser,
		db: database,
	}
}

func (reader Reader)Parse(feedUrl string, jobs chan<- string) error {
	feed, err := reader.Parser.ParseURL(feedUrl)
	if err != nil { return err }

	channelName := feed.Title

	lastURL, err := reader.db.GetLastItemURL(feedUrl)
	if err != nil { return err }

	var firstItemURL string

	for _, item := range feed.Items {
		url := item.Link
		if strings.Contains(url, "shorts") {
			continue
		}

		// Track the first item's URL (newest item)
		if firstItemURL == "" {
			firstItemURL = url
		}

		// If this is the item we last notified about, stop processing
		if lastURL != "" && url == lastURL {
			break
		}

		title := item.Title
		jobs <- formatMessage(channelName, url, title)
	}

	err = reader.db.SaveLastItemURL(feedUrl, firstItemURL)
	if err != nil { return err }
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
