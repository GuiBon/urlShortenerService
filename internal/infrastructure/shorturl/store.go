package shorturl

import (
	"context"
	"errors"
	"time"
	"urlShortenerService/domain"
)

var (
	// ErrNotFound is the error when a slug is not found within the database
	ErrNotFound error = errors.New("url not found")
)

// Store represents operations on shorturl Store
type Store interface {
	// DeleteExpired deletes the slug / URL couples that are expired
	DeleteExpired(ctx context.Context, timeToExpire time.Duration) (int, error)
	// Get retrieves the URL associated to a specific slug
	Get(ctx context.Context, slug string) (domain.URLMapping, error)
	// Set stores the slug and the URL associated
	Set(ctx context.Context, shortURL domain.URLMapping) error
}
