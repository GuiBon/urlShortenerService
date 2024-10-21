package shorturl

import (
	context "context"
	"sync"
	"time"
	"urlShortenerService/domain"
)

// CacheStore represents a inmemory cache store
type CacheStore struct {
	persistentStore Store
	cacheStore      sync.Map
}

// NewCacheStore creates a cache store
func NewCacheStore(persistentStore Store) *CacheStore {
	return &CacheStore{
		persistentStore: persistentStore,
		cacheStore:      sync.Map{},
	}
}

// Set implements Store interface
func (s *CacheStore) Set(ctx context.Context, shortURL domain.URLMapping) error {
	err := s.persistentStore.Set(ctx, shortURL)
	if err != nil {
		return err
	}

	s.cacheStore.Store(shortURL.Slug, shortURL.OriginalURL)
	return nil
}

// Get implements Store interface
func (s *CacheStore) Get(ctx context.Context, slug string) (domain.URLMapping, error) {
	originalURL, exists := s.cacheStore.Load(slug)
	if exists {
		return domain.URLMapping{
			Slug:        slug,
			OriginalURL: originalURL.(string),
		}, nil
	}
	return s.persistentStore.Get(ctx, slug)
}

func (s *CacheStore) DeleteExpired(ctx context.Context, timeToExpire time.Duration) ([]string, error) {
	slugsDeleted, err := s.persistentStore.DeleteExpired(ctx, timeToExpire)
	if err != nil {
		return nil, err
	}
	for _, slug := range slugsDeleted {
		s.cacheStore.Delete(slug)
	}
	return slugsDeleted, nil
}
