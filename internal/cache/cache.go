package cache

import (
	"time"
	"os"
)

const CACHE_FILE string = "cache"

func ReadLastCheck() (*time.Time, error) {
	content, err := os.ReadFile(CACHE_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	res, err := time.Parse(time.RFC3339, string(content))
	return &res, err
}

func WriteLastCheck() error {
	return os.WriteFile(CACHE_FILE, []byte(time.Now().Format(time.RFC3339)), 0644)
}
