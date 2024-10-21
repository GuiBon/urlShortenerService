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
	// GetURL retrieves the statistic for a single URL
	GetURL(ctx context.Context, url string) (domain.URLStatistic, error)
	// GetTopURLs retrieves top statistic of the choosen type for the URLs
	GetTopURLs(ctx context.Context, statType StatisticType, limitOveride int64) ([]domain.URLStatistic, error)
	// SetURL stores the statistic of the choosen type for the associated URL
	SetURL(ctx context.Context, url string, statType StatisticType) error
}
