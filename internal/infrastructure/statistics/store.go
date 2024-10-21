package statistics

import (
	"context"
	"urlShortenerService/domain"
)

type StatisticType string

var (
	StatisticTypeShortened StatisticType = "urls-shortened"
	StatisticTypeAccessed  StatisticType = "urls-accessed"
)

// Store represents operations on statistics Store
type Store interface {
	// GetOne retrieves the statistic for a single URL
	GetOne(ctx context.Context, url string) (domain.URLStatistic, error)
	// GetTop retrieves top statistic of the choosen type for the URLs
	GetTop(ctx context.Context, statType StatisticType, limitOveride int64) ([]domain.URLStatistic, error)
	// Set stores the statistic of the choosen type for the associated URL
	Set(ctx context.Context, url string, statType StatisticType) error
}
