package usecase

import (
	"context"
	"fmt"
	"urlShortenerService/domain"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/infrastructure/shorturl"
	"urlShortenerService/internal/infrastructure/statistics"

	"github.com/golang/glog"
)

// CreateShortenURLCmd represents the function signature of the command that create a shorten URL
type CreateShortenURLCmd func(ctx context.Context, urlToShorten string) (string, error)

// createShortenURL creates, stores and returns a shorten URL
func createShortenURL(baseURL string, urlSanitizerCmd command.URLSanitizerCmd, slugGeneratorCmd command.SlugGeneratorCmd,
	shortURLStore shorturl.Store, statisticsStore statistics.Store) CreateShortenURLCmd {
	return func(ctx context.Context, urlToShorten string) (string, error) {
		// Sanitize and validate URL
		sanitizedURLToShorten, err := urlSanitizerCmd(urlToShorten)
		if err != nil {
			return "", err
		}

		// Shorten URL
		slug := slugGeneratorCmd(sanitizedURLToShorten)

		// Save URL
		err = shortURLStore.Set(ctx, domain.URLMapping{
			Slug:        slug,
			OriginalURL: sanitizedURLToShorten,
		})
		if err != nil {
			return "", err
		}

		// Update statistics
		go func(url string) {
			err := statisticsStore.SetURL(context.Background(), url, statistics.StatisticTypeShortened)
			if err != nil {
				glog.Errorf("failed to set [%s] statistics for [%s]: %w", statistics.StatisticTypeShortened, url, err)
			}
		}(sanitizedURLToShorten)

		return fmt.Sprintf("%s/%s", baseURL, slug), nil
	}
}

// CreateShortenURLCmdBuilder builds the command that will create a shorten URL
func CreateShortenURLCmdBuilder(baseURL string, urlSanitizerCmd command.URLSanitizerCmd, slugGeneratorCmd command.SlugGeneratorCmd,
	shortURLStore shorturl.Store, statisticsStore statistics.Store) CreateShortenURLCmd {
	return createShortenURL(baseURL, urlSanitizerCmd, slugGeneratorCmd, shortURLStore, statisticsStore)
}
