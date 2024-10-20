package domain

import "time"

// URLMapping represents an URL mapping data between a short URL and its original form
type URLMapping struct {
	Slug        string    `db:"slug"`
	OriginalURL string    `db:"original_url"`
	InsertedAt  time.Time `db:"inserted_at"`
}
