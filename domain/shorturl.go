package domain

import "time"

// ShortURL represents a shortened URL data
type ShortURL struct {
	Slug       string    `db:"slug"`
	URL        string    `db:"url"`
	InsertedAt time.Time `db:"inserted_at"`
}
