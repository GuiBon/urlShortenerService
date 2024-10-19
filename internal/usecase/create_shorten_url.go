package usecase

import (
	"context"
	"fmt"
	"urlShortenerService/domain"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/infrastructure/shorturl"
)

// CreateShortenURLCmd represents the function signature of the command that create a shorten URL
type CreateShortenURLCmd func(ctx context.Context, urlToShorten string) (string, error)

// createShortenURL creates, stores and returns a shorten URL
func createShortenURL(baseURL string, urlSanitizerCmd command.URLSanitizerCmd,
	slugGeneratorCmd command.SlugGeneratorCmd, shortURLStore shorturl.Store) CreateShortenURLCmd {
	return func(ctx context.Context, urlToShorten string) (string, error) {
		// Sanitize and validate URL
		sanitizedURLToShorten, err := urlSanitizerCmd(urlToShorten)
		if err != nil {
			return "", err
		}

		// Shorten URL
		slug := slugGeneratorCmd(sanitizedURLToShorten)

		// Save URL
		err = shortURLStore.Set(ctx, domain.ShortURL{
			Slug: slug,
			URL:  sanitizedURLToShorten,
		})

		return fmt.Sprintf("%s/%s", baseURL, slug), err
	}
}

// CreateShortenURLCmdBuilder builds the command that will create a shorten URL
func CreateShortenURLCmdBuilder(baseURL string, urlSanitizerCmd command.URLSanitizerCmd,
	slugGeneratorCmd command.SlugGeneratorCmd, shortURLStore shorturl.Store) CreateShortenURLCmd {
	return createShortenURL(baseURL, urlSanitizerCmd, slugGeneratorCmd, shortURLStore)
}
