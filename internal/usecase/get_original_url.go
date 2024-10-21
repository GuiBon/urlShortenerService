package usecase

import (
	"context"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/infrastructure/shorturl"
	"urlShortenerService/internal/infrastructure/statistics"

	"github.com/golang/glog"
)

// GetOriginalURLCmd represents the function signature of the command that retrieves an original URL given a slug
type GetOriginalURLCmd func(ctx context.Context, shortURL string) (string, error)

// getOriginalURL retrieves an original URL given a slug
func getOriginalURL(slugValidatorCmd command.SlugValidatorCmd, shortURLStore shorturl.Store, statisticsStore statistics.Store) GetOriginalURLCmd {
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

		// Update statistics
		go func(url string) {
			err := statisticsStore.SetURL(context.Background(), url, statistics.StatisticTypeAccessed)
			if err != nil {
				glog.Errorf("failed to set [%s] statistics for [%s]: %w", statistics.StatisticTypeAccessed, url, err)
			}
		}(urlMapping.OriginalURL)

		return urlMapping.OriginalURL, nil
	}
}

// GetOriginalURLCmdBuilder builds the command that will retrieves an original URL
func GetOriginalURLCmdBuilder(slugValidatorCmd command.SlugValidatorCmd, shortURLStore shorturl.Store, statisticsStore statistics.Store) GetOriginalURLCmd {
	return getOriginalURL(slugValidatorCmd, shortURLStore, statisticsStore)
}
