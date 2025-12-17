package config

import (
	"os"
	"time"
	"fmt"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	FeedsPath string
	MaxItemAge time.Duration

	DatabasePath string

	TelegramChatID string
	TelegramToken string

	HTTPTimeout time.Duration
}

type Feed struct {
	URL  string `yaml:"url"`
	Type string `yaml:"type"`
}

type FeedsConfig struct {
	Feeds []Feed
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		FeedsPath: getEnv("FEEDS_PATH", "config/feeds.yml"),
		MaxItemAge: 24 * time.Hour,
		DatabasePath: getEnv("DATABASE_PATH", "cache.db"),
		TelegramChatID: os.Getenv("TELEGRAM_CHAT_ID"),
		TelegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		HTTPTimeout: 30 * time.Second,
	}

	if err := cfg.validate(); err != nil { return nil, err }

	return cfg, nil
}

func (cfg *Config)LoadFeeds() ([]Feed, error) {
	data, err := os.ReadFile(cfg.FeedsPath)
	if err != nil { return nil, fmt.Errorf("failed to read feeds config file: %v", err) }

	var feedsConfig FeedsConfig
	if err := yaml.Unmarshal(data, &feedsConfig); err != nil {
		return nil, fmt.Errorf("failed to parse feeds config file: %v", err)
	}

	return feedsConfig.Feeds, nil
}

func (cfg *Config)validate() error {
	if cfg.TelegramChatID == "" {
		return fmt.Errorf("TELEGRAM_CHAT_ID is required")
	}
	if cfg.TelegramToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
