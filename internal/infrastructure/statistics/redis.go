package statistics

import (
	"context"
	"fmt"
	"sync"
	"urlShortenerService/domain"
	"urlShortenerService/internal/infrastructure/config"

	"github.com/go-redis/redis/v8"
)

// RedisStore represents a redis store thread proof
type RedisStore struct {
	client     *redis.Client
	mutex      sync.RWMutex
	maxResults int64
}

// NewRedisStore connects to a redis and return it inside a RedisStore
func NewRedisStore(cfg config.RedisConfig) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.ToAddr(),
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisStore{
		client:     client,
		mutex:      sync.RWMutex{},
		maxResults: int64(cfg.MaxResults),
	}, nil
}

// GetOne implements the Store interface
func (s *RedisStore) GetOne(ctx context.Context, url string) (domain.URLStatistic, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	shortened, err := s.client.ZScore(ctx, string(StatisticTypeShortened), url).Result()
	if err != nil && err != redis.Nil {
		return domain.URLStatistic{}, fmt.Errorf("failed to get [%s] stats for URL [%s]: %w", StatisticTypeShortened, url, err)
	}

	accessed, err := s.client.ZScore(ctx, string(StatisticTypeAccessed), url).Result()
	if err != nil && err != redis.Nil {
		return domain.URLStatistic{}, fmt.Errorf("failed to get [%s] stats for URL [%s]: %w", StatisticTypeAccessed, url, err)
	}

	return domain.URLStatistic{
		URL:              url,
		ShortenedCounter: int(shortened),
		AccessedCounter:  int(accessed),
	}, nil
}

// GetTop implements the Store interface
func (s *RedisStore) GetTop(ctx context.Context, statType StatisticType, limitOveride int64) ([]domain.URLStatistic, error) {
	var limit = s.maxResults
	if limitOveride != 0 {
		limit = limitOveride
	}

	s.mutex.RLock()
	zSlice, err := s.client.ZRevRangeWithScores(ctx, string(statType), 0, limit-1).Result()
	if err != nil {
		s.mutex.RUnlock()
		return nil, fmt.Errorf("failed to get [%s] top stats: %w", statType, err)
	}
	s.mutex.RUnlock()

	var stats []domain.URLStatistic
	for _, z := range zSlice {
		switch statType {
		case StatisticTypeShortened:
			stats = append(stats, domain.URLStatistic{
				URL:              z.Member.(string),
				ShortenedCounter: int(z.Score),
			})
		case StatisticTypeAccessed:
			stats = append(stats, domain.URLStatistic{
				URL:             z.Member.(string),
				AccessedCounter: int(z.Score),
			})
		}
	}

	return stats, nil
}

// Set implements the Store interface
func (s *RedisStore) Set(ctx context.Context, url string, statType StatisticType) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, err := s.client.ZIncrBy(ctx, string(statType), 1, url).Result()
	if err != nil {
		return fmt.Errorf("failed to set [%s] stat for URL [%s]: %w", statType, url, err)
	}

	return nil
}
