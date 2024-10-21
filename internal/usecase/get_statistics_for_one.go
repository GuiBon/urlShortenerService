package usecase

import (
	"context"
	"urlShortenerService/domain"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/infrastructure/statistics"
)

// GetStatisticsForOneCmd represents the function signature of the command that retrieves statistics for a given URL
type GetStatisticsForOneCmd func(ctx context.Context, url string) (domain.URLStatistic, error)

// getStatisticsForOne retrieves statistics for a given URL
func getStatisticsForOne(urlSanitizerCmd command.URLSanitizerCmd, statisticsStore statistics.Store) GetStatisticsForOneCmd {
	return func(ctx context.Context, url string) (domain.URLStatistic, error) {
		// Sanitize and validate URL
		sanitizedURL, err := urlSanitizerCmd(url)
		if err != nil {
			return domain.URLStatistic{}, err
		}

		// Retrieves statistics
		return statisticsStore.GetOne(ctx, sanitizedURL)
	}
}

// GetStatisticsForOneCmdBuilder builds the command that will retrieves statistics
func GetStatisticsForOneCmdBuilder(urlSanitizerCmd command.URLSanitizerCmd, statisticsStore statistics.Store) GetStatisticsForOneCmd {
	return getStatisticsForOne(urlSanitizerCmd, statisticsStore)
}
