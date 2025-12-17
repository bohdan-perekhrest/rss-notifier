package main

import (
	"log"

	"rss-notifier/config"
	"rss-notifier/internal/app"
	"rss-notifier/internal/cache"
	"rss-notifier/internal/notifier"
	"rss-notifier/internal/reader"
	"rss-notifier/pkg/httpclient"
)

func main() {
	cfg, err := config.Load()
	if err != nil { log.Fatalf("Failed to load config: %v", err) }

	db, err := cache.NewSQLiteCache(cfg.DatabasePath)
	if err != nil { log.Fatalf("Failed to initialize cache: %v", err) }
	defer db.Close()

	httpClient := httpclient.New(cfg.HTTPTimeout)
	feedReader := reader.New(httpClient, db, cfg.MaxItemAge)
	telegramNotifier := notifier.NewTelegram(httpClient, cfg.TelegramChatID, cfg.TelegramToken)

	svc := app.NewService(feedReader, telegramNotifier, db, cfg)

	if err := svc.Run(); err != nil { log.Fatalf("Service failed: %v", err) }

}
