package main

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"rss-notifier/internal/cache"
	"rss-notifier/internal/reader"
	"sync"
	"time"
)

func processFeed(url string, lastCheck *time.Time, wg *sync.WaitGroup) {
	defer wg.Done()

	reader := reader.NewReader()
	if err := reader.Parse(url, lastCheck); err != nil {
		log.Printf("Failed to parse feed %s: %v\n", url, err)
	}
}

func main() {
	yfile, err := os.ReadFile("config/feeds.yml")
	if err != nil { log.Fatal(err) }

	var data map[string][]map[string]string
	err = yaml.Unmarshal(yfile, &data)
	if err != nil { log.Fatal(err) }

	lastCheck, err := cache.ReadLastCheck()
	if err != nil { log.Printf("Warning: could not read cache: %v\n", err) }

	var wg sync.WaitGroup

	for _, feed := range data["feeds"] {
		url := feed["url"]
		if feed["type"] == "youtube" {
			url = "https://www.youtube.com/feeds/videos.xml?channel_id=" + feed["url"]
		}

		wg.Add(1)
		go processFeed(url, lastCheck, &wg)
	}
	wg.Wait()

	if err := cache.WriteLastCheck(); err != nil { log.Fatal(err) }
}
