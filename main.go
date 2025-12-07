package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"rss-notifier/internal/cache"
	"rss-notifier/internal/reader"
	"rss-notifier/internal/telegram"

	"gopkg.in/yaml.v3"
)

var (
	client = createHTTPClient()
	Reader = reader.NewReader(client)
	Sender = telegram.NewSender(client)
)

func createHTTPClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
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

	messagesJobs := make(chan string, 5)
	var workerWg sync.WaitGroup
	var feedWg sync.WaitGroup

	for i := 1; i <= 4; i++ {
		workerWg.Add(1)
		go func(jobs <-chan string) {
			defer workerWg.Done()

			for text := range jobs {
				if err := Sender.SendMessage(text); err != nil { log.Printf("Failed to send message: %v\n", err) }
			}
		}(messagesJobs)
	}

	for _, feed := range data["feeds"] {
		url := feed["url"]
		if feed["type"] == "youtube" {
			url = "https://www.youtube.com/feeds/videos.xml?channel_id=" + feed["url"]
		}

		feedWg.Add(1)
		go func(url string) {
			defer feedWg.Done()

			if err := Reader.Parse(url, lastCheck, messagesJobs); err != nil {
				log.Printf("Failed to parse feed %s: %v\n", url, err)
			}
		}(url)
	}

	feedWg.Wait()
	close(messagesJobs)
	workerWg.Wait()

	if err := cache.WriteLastCheck(); err != nil { log.Fatal(err) }
}
