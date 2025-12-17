package feed

import "time"

type Feed struct {
	URL      string
	Type     string
	Provider Provider
}

type Item struct {
	URL       string
	Title     string
	Published time.Time
}

type Response struct {
	FeedURL 	    string
	ChannelName   string
	Items 			  []Item
	NewestItemURL string
}
