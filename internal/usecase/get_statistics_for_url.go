package usecase

import (
	"context"
	"urlShortenerService/domain"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/infrastructure/statistics"
)

// GetStatisticsForURLCmd represents the function signature of the command that retrieves statistics for a given URL
type GetStatisticsForURLCmd func(ctx context.Context, url string) (domain.URLStatistic, error)

// getStatisticsForURL retrieves statistics for a given URL
func getStatisticsForURL(urlSanitizerCmd command.URLSanitizerCmd, statisticsStore statistics.Store) GetStatisticsForURLCmd {
	return func(ctx context.Context, url string) (domain.URLStatistic, error) {
		// Sanitize and validate URL
		sanitizedURL, err := urlSanitizerCmd(url)
		if err != nil {
			return domain.URLStatistic{}, err
		}

		// Retrieves statistics
		return statisticsStore.GetURL(ctx, sanitizedURL)
	}
}

// GetStatisticsForURLCmdBuilder builds the command that will retrieves statistics
func GetStatisticsForURLCmdBuilder(urlSanitizerCmd command.URLSanitizerCmd, statisticsStore statistics.Store) GetStatisticsForURLCmd {
	return getStatisticsForURL(urlSanitizerCmd, statisticsStore)
}
