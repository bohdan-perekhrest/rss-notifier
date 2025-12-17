package notifier

import (
	"html"
	"strings"

	"rss-notifier/internal/feed"
)

func FormatFeeds(responses []feed.Response) string {
	var builder strings.Builder

	for _, resp := range responses {
		builder.WriteString("FROM: <b>")
		builder.WriteString(html.EscapeString(resp.ChannelName))
		builder.WriteString("</b>\n")

		for _, item := range resp.Items {
			builder.WriteString("<a href=\"")
			builder.WriteString(item.URL)
			builder.WriteString("\">")
			builder.WriteString(html.EscapeString(item.Title))
			builder.WriteString("</a>\n")
		}

		builder.WriteString("\n")
	}

	return builder.String()
}
