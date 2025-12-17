package cache

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type SQLiteCache struct {
	db *sql.DB
}

func NewSQLiteCache(path string) (*SQLiteCache, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if _, err := db.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set busy_timeout: %w", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set journal_mode: %w", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS feed_cache (
		feed_url TEXT PRIMARY KEY,
		last_item_url TEXT NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &SQLiteCache{db: db}, nil
}

func (c *SQLiteCache) GetLastItemURL(url string) (string, error) {
	var result string
	query := "SELECT last_item_url FROM feed_cache WHERE feed_url = ?"
	err := c.db.QueryRow(query, url).Scan(&result)

	if err == sql.ErrNoRows {
		return "", nil // Not an error, just no record
	}

	if err != nil {
		return "", fmt.Errorf("failed to query cache: %w", err)
	}

	return result, nil
}

func (c *SQLiteCache) SaveLastItemURL(urls map[string]string) error {
	if len(urls) == 0 {
		return nil
	}

	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
	INSERT INTO feed_cache (feed_url, last_item_url, updated_at)
	VALUES (?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(feed_url)
	DO UPDATE SET last_item_url = ?, updated_at = CURRENT_TIMESTAMP`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for feedURL, itemURL := range urls {
		if _, err := stmt.Exec(feedURL, itemURL, itemURL); err != nil {
			return fmt.Errorf("failed to save %s: %w", feedURL, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (c *SQLiteCache) Close() error {
	return c.db.Close()
}
