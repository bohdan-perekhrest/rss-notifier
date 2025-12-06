package main

import (
	"log"
	"os"
	"time"

	"rss-notifier/internal/cache"
	"rss-notifier/internal/reader"

	"gopkg.in/yaml.v3"
)

func processFeed(url string, lastCheck *time.Time, done chan<- bool) {
	reader := reader.NewReader()
	if err := reader.Parse(url, lastCheck); err != nil {
		log.Printf("Failed to parse feed %s: %v\n", url, err)
	}

	done <- true
}

func main() {
	yfile, err := os.ReadFile("config/feeds.yml")
	if err != nil { log.Fatal(err) }

	var data map[string][]map[string]string
	err = yaml.Unmarshal(yfile, &data)
	if err != nil { log.Fatal(err) }

	lastCheck, err := cache.ReadLastCheck()
	if err != nil { log.Printf("Warning: could not read cache: %v\n", err) }

	done := make(chan bool, len(data["feeds"]))

	for _, feed := range data["feeds"] {
		url := feed["url"]
		if feed["type"] == "youtube" {
			url = "https://www.youtube.com/feeds/videos.xml?channel_id=" + feed["url"]
		}

		go processFeed(url, lastCheck, done)
	}

	for i := 0; i < len(data["feeds"]); i++ {
		<-done
	}

	if err := cache.WriteLastCheck(); err != nil { log.Fatal(err) }
}
