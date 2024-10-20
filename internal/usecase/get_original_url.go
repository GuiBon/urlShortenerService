package usecase

import (
	"context"
	"urlShortenerService/internal/infrastructure/shorturl"
)

// GetOriginalURLCmd represents the function signature of the command that retrieves an original URL given a slug
type GetOriginalURLCmd func(ctx context.Context, shortURL string) (string, error)

// getOriginalURL retrives an original URL given a slug
func getOriginalURL(shortURLStore shorturl.Store) GetOriginalURLCmd {
	return func(ctx context.Context, slug string) (string, error) {
		// Retrieves URL
		urlMapping, err := shortURLStore.Get(ctx, slug)
		if err != nil {
			return "", err
		}

		return urlMapping.OriginalURL, nil
	}
}

// GetOriginalURLCmdBuilder builds the command that will retrieves an original URL
func GetOriginalURLCmdBuilder(shortURLStore shorturl.Store) GetOriginalURLCmd {
	return getOriginalURL(shortURLStore)
}
