package database

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

const DATABASE_FILE string = "cache.db"

type Database struct {
	db *sql.DB
}

func InitDB() (*Database, error) {
	db, err := sql.Open("sqlite", DATABASE_FILE)
	if err != nil {
		return nil, err
	}

	// Set busy timeout: wait up to 5 seconds if database is locked
	_, err = db.Exec("PRAGMA busy_timeout = 5000")
	if err != nil {
		db.Close()
		return nil, err
	}

	// Enable WAL mode for better concurrent access
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		db.Close()
		return nil, err
	}

	// Create table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS feed_cache (
		feed_url TEXT PRIMARY KEY,
		last_item_url TEXT NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Database{db: db}, nil
}

func (db Database)GetLastItemURL(feedURL string) (string, error) {
	var lastItemURL string
	query := "SELECT last_item_url FROM feed_cache WHERE feed_url = ?"
	err := db.db.QueryRow(query, feedURL).Scan(&lastItemURL)

	if err == sql.ErrNoRows {
		// No record found, return empty string (not an error)
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return lastItemURL, nil
}

func (db Database)SaveLastItemURL(feedURL, lastItemURL string) error {
	query := `
	INSERT INTO feed_cache (feed_url, last_item_url, updated_at)
	VALUES (?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(feed_url)
	DO UPDATE SET last_item_url = ?, updated_at = CURRENT_TIMESTAMP`

	_, err := db.db.Exec(query, feedURL, lastItemURL, lastItemURL)
	return err
}

func (db Database)CloseDB() error {
	return db.db.Close()
}
