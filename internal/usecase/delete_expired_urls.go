package usecase

import (
	"context"
	"time"
	"urlShortenerService/internal/infrastructure/shorturl"
)

// DeleteExpiredURLsCmd represents the function signature of the command that deletes expired URLs
type DeleteExpiredURLsCmd func(ctx context.Context) (int, error)

// deleteExpiredURLs deletes URLs that have expired
func deleteExpiredURLs(timeToExpire time.Duration, shortURLStore shorturl.Store) DeleteExpiredURLsCmd {
	return func(ctx context.Context) (int, error) {
		// Deletes expired URL
		return shortURLStore.DeleteExpired(ctx, timeToExpire)
	}
}

// DeleteExpiredURLsCmdBuilder builds the command that will deletes expired URLs
func DeleteExpiredURLsCmdBuilder(timeToExpire time.Duration, shortURLStore shorturl.Store) DeleteExpiredURLsCmd {
	return deleteExpiredURLs(timeToExpire, shortURLStore)
}
