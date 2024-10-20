package usecase

import (
	"context"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/infrastructure/shorturl"
)

// GetOriginalURLCmd represents the function signature of the command that retrieves an original URL given a slug
type GetOriginalURLCmd func(ctx context.Context, shortURL string) (string, error)

// getOriginalURL retrieves an original URL given a slug
func getOriginalURL(slugValidatorCmd command.SlugValidatorCmd, shortURLStore shorturl.Store) GetOriginalURLCmd {
	return func(ctx context.Context, slug string) (string, error) {
		// Ensure slug validity to avoid useless query to store
		err := slugValidatorCmd(slug)
		if err != nil {
			return "", err
		}

		// Retrieves URL
		urlMapping, err := shortURLStore.Get(ctx, slug)
		if err != nil {
			return "", err
		}

		return urlMapping.OriginalURL, nil
	}
}

// GetOriginalURLCmdBuilder builds the command that will retrieves an original URL
func GetOriginalURLCmdBuilder(slugValidatorCmd command.SlugValidatorCmd, shortURLStore shorturl.Store) GetOriginalURLCmd {
	return getOriginalURL(slugValidatorCmd, shortURLStore)
}
