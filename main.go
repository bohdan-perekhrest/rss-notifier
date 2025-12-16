package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"rss-notifier/internal/database"
	"rss-notifier/internal/reader"
	"rss-notifier/internal/telegram"

	"gopkg.in/yaml.v3"
)

func createHTTPClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     10 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

func readFeedsConfig() ([]map[string]string, error) {
	yfile, err := os.ReadFile("config/feeds.yml")
	if err != nil { return nil, err }

	var data map[string][]map[string]string
	err = yaml.Unmarshal(yfile, &data)
	if err != nil { return nil, err }

	return data["feeds"], nil
}

func main() {
	// Load feeds
	feeds, err := readFeedsConfig()
	if err != nil { log.Fatal(err) }

	// Initialize Database
	db, err := database.InitDB()
	if err != nil { log.Fatal(err) }
	defer db.CloseDB()

	client := createHTTPClient()
	parser := reader.NewReader(client, db)
	sender := telegram.NewSender(client)

	messagesJobs := make(chan string, 5)
	var workerWg sync.WaitGroup
	var feedWg sync.WaitGroup

	for i := 1; i <= 4; i++ {
		workerWg.Add(1)
		go func(jobs <-chan string) {
			defer workerWg.Done()

			for text := range jobs {
				if err := sender.SendMessage(text); err != nil { log.Printf("Failed to send message: %v\n", err) }
			}
		}(messagesJobs)
	}

	for _, feed := range feeds {
		url := feed["url"]
		if feed["type"] == "youtube" {
			url = "https://www.youtube.com/feeds/videos.xml?channel_id=" + feed["url"]
		}

		feedWg.Add(1)
		go func(url string) {
			defer feedWg.Done()

			err := parser.Parse(url, messagesJobs)
			if err != nil {
				log.Printf("Failed to parse feed %s: %v\n", url, err)
				return
			}
		}(url)
	}

	feedWg.Wait()
	close(messagesJobs)
	workerWg.Wait()
}
