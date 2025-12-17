package feed

import (
	"time"
	"github.com/mmcdole/gofeed"
)

type Filter struct {
	maxAge time.Duration
}

func NewFilter(maxAge time.Duration) *Filter {
	return &Filter{maxAge: maxAge}
}

func (f *Filter) ShouldInclude(item *gofeed.Item) bool {
	if item.PublishedParsed == nil {
		return false
	}

	cutoff := time.Now().Add(-f.maxAge)
	return item.PublishedParsed.After(cutoff)
}
