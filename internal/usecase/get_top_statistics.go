package usecase

import (
	"context"
	"urlShortenerService/domain"
	"urlShortenerService/internal/infrastructure/statistics"
)

// GetTopStatisticsCmd represents the function signature of the command that retrieves top statistics for a given statistic type
type GetTopStatisticsCmd func(ctx context.Context, statType statistics.StatisticType, limitOveride int64) ([]domain.URLStatistic, error)

// getTopStatistics retrieves top statistics for a given statistic type
func getTopStatistics(statisticsStore statistics.Store) GetTopStatisticsCmd {
	return func(ctx context.Context, statType statistics.StatisticType, limitOveride int64) ([]domain.URLStatistic, error) {
		return statisticsStore.GetTopURLs(ctx, statType, limitOveride)
	}
}

// GetTopStatisticsCmdBuilder builds the command that will retrieves top statistics
func GetTopStatisticsCmdBuilder(statisticsStore statistics.Store) GetTopStatisticsCmd {
	return getTopStatistics(statisticsStore)
}
