package shorturl

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
	"urlShortenerService/domain"
	"urlShortenerService/internal/infrastructure/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCacheStore(t *testing.T) {
	os.Setenv("env", "test")
	defer os.Unsetenv("env")
	conf, err := config.Load()
	require.NoError(t, err)
	persistentStore, err := NewPSQLStore(conf.Database)
	require.NoError(t, err)

	store := NewCacheStore(persistentStore)

	RunStoreTests(t, store)
}

func TestCacheSet(t *testing.T) {
	shortURL := domain.URLMapping{
		Slug:        "example",
		OriginalURL: "https://example.com",
	}
	t.Run("nominal", func(t *testing.T) {
		// Given
		persitentMockStore := NewMock(t)
		persitentMockStore.On("Set", mock.Anything, shortURL).Return(nil)
		store := NewCacheStore(persitentMockStore)

		// When
		err := store.Set(context.Background(), shortURL)
		require.NoError(t, err)

		// Then
		originalURL, exists := store.cacheStore.Load(shortURL.Slug)
		assert.True(t, exists)
		assert.Equal(t, shortURL.OriginalURL, originalURL.(string))
	})
	t.Run("with persistent store failed", func(t *testing.T) {
		// Given
		persitentMockStore := NewMock(t)
		persitentMockStore.On("Set", mock.Anything, shortURL).Return(assert.AnError)
		store := NewCacheStore(persitentMockStore)

		// When
		err := store.Set(context.Background(), shortURL)

		// Then
		assert.ErrorIs(t, err, assert.AnError)
		originalURL, exists := store.cacheStore.Load(shortURL.Slug)
		assert.False(t, exists)
		assert.Empty(t, originalURL)
	})
}

func TestCacheGet(t *testing.T) {
	slug := "jV6gHv0o"
	shortURL := domain.URLMapping{
		Slug:        slug,
		OriginalURL: "https://example.com",
	}
	t.Run("found in cache", func(t *testing.T) {
		// Given
		persitentMockStore := NewMock(t)
		store := NewCacheStore(persitentMockStore)
		store.cacheStore.Store(slug, shortURL.OriginalURL)

		// When
		urlMapping, err := store.Get(context.Background(), slug)
		require.NoError(t, err)

		// Then
		assert.Equal(t, shortURL, urlMapping)
	})
	t.Run("found in persistent store", func(t *testing.T) {
		// Given
		persitentMockStore := NewMock(t)
		persitentMockStore.On("Get", mock.Anything, slug).Return(shortURL, nil)
		store := NewCacheStore(persitentMockStore)

		// When
		urlMapping, err := store.Get(context.Background(), slug)
		require.NoError(t, err)

		// Then
		assert.Equal(t, shortURL, urlMapping)
	})
	t.Run("persistent store errored", func(t *testing.T) {
		// Given
		persitentMockStore := NewMock(t)
		persitentMockStore.On("Get", mock.Anything, slug).Return(domain.URLMapping{}, assert.AnError)
		store := NewCacheStore(persitentMockStore)

		// When
		urlMapping, err := store.Get(context.Background(), slug)

		// Then
		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, urlMapping)
	})
}

func TestCacheDeleteExpired(t *testing.T) {
	timeToExpire := 10 * time.Minute
	slugsToDelete := []string{"2zv8a2Im", "1eJSWjFM", "UsIJeS1D", "K11q8dTj", "Sd7k2eDU"}
	t.Run("nominal", func(t *testing.T) {
		// Given
		persitentMockStore := NewMock(t)
		persitentMockStore.On("DeleteExpired", mock.Anything, timeToExpire).Return(slugsToDelete, nil)
		store := NewCacheStore(persitentMockStore)
		for i, slug := range slugsToDelete {
			store.cacheStore.Store(slug, fmt.Sprintf("https://example.com/%d", i))
		}

		// When
		slugsDeleted, err := store.DeleteExpired(context.Background(), timeToExpire)
		require.NoError(t, err)

		// Then
		assert.Equal(t, slugsToDelete, slugsDeleted)
		for _, slugDeleted := range slugsDeleted {
			originalURL, exists := store.cacheStore.Load(slugDeleted)
			assert.False(t, exists)
			assert.Empty(t, originalURL)
		}
	})
	t.Run("persistent store errored", func(t *testing.T) {
		// Given
		persitentMockStore := NewMock(t)
		persitentMockStore.On("DeleteExpired", mock.Anything, timeToExpire).Return([]string{}, assert.AnError)
		store := NewCacheStore(persitentMockStore)
		for i, slug := range slugsToDelete {
			store.cacheStore.Store(slug, fmt.Sprintf("https://example.com/%d", i))
		}

		// When
		slugsDeleted, err := store.DeleteExpired(context.Background(), timeToExpire)

		// Then
		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, slugsDeleted)
		for _, slugDeleted := range slugsDeleted {
			originalURL, exists := store.cacheStore.Load(slugDeleted)
			assert.True(t, exists)
			assert.NotEmpty(t, originalURL)
		}
	})
}
