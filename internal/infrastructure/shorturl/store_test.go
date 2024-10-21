package shorturl

import (
	"context"
	"testing"
	time "time"
	"urlShortenerService/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StoreTestSuite struct {
	Store
}

func RunStoreTests(t *testing.T, store Store) {
	suite := &StoreTestSuite{Store: store}

	t.Run("TestSet", suite.TestSet)
	t.Run("TestGet", suite.TestGet)
	t.Run("TestSetDuplicateSlug", suite.TestSetDuplicateSlug)
	t.Run("TestDeleteExpired", suite.TestDeleteExpired)
}

func (suite *StoreTestSuite) TestSet(t *testing.T) {
	// Given
	ctx := context.Background()
	shortURL := domain.URLMapping{
		Slug:        "example",
		OriginalURL: "https://example.com",
	}

	// When
	err := suite.Store.Set(ctx, shortURL)
	require.NoError(t, err)

	// Then
	retrievedURL, err := suite.Store.Get(ctx, shortURL.Slug)
	require.NoError(t, err)
	assert.Equal(t, shortURL.Slug, retrievedURL.Slug)
	assert.Equal(t, shortURL.OriginalURL, retrievedURL.OriginalURL)
	if _, ok := suite.Store.(*CacheStore); !ok { // This assertion shouldn't be tested for cache
		assert.NotEmpty(t, retrievedURL.InsertedAt)
	}
}

func (suite *StoreTestSuite) TestGet(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		// Given
		ctx := context.Background()
		shortURL := domain.URLMapping{
			Slug:        "example",
			OriginalURL: "https://example.com",
		}
		err := suite.Store.Set(ctx, shortURL)
		require.NoError(t, err)

		// When
		retrievedURL, err := suite.Store.Get(ctx, shortURL.Slug)
		require.NoError(t, err)

		// Then
		assert.Equal(t, shortURL.Slug, retrievedURL.Slug)
		assert.Equal(t, shortURL.OriginalURL, retrievedURL.OriginalURL)
		if _, ok := suite.Store.(*CacheStore); !ok { // This assertion shouldn't be tested for cache
			assert.NotEmpty(t, retrievedURL.InsertedAt)
		}
	})
	t.Run("slug not found", func(t *testing.T) {
		// Given
		ctx := context.Background()

		// When
		retrievedURL, err := suite.Store.Get(ctx, "unknown-slug")

		// Then
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Empty(t, retrievedURL)
	})
}

func (suite *StoreTestSuite) TestSetDuplicateSlug(t *testing.T) {
	// Given
	ctx := context.Background()
	shortURL1 := domain.URLMapping{
		Slug:        "duplicate",
		OriginalURL: "https://example.com/duplicate",
		InsertedAt:  time.Now().Add(-2 * time.Hour).UTC(),
	}
	shortURL2 := domain.URLMapping{
		Slug:        "duplicate",
		OriginalURL: "https://example.com/duplicate",
		InsertedAt:  time.Now().Add(-1 * time.Hour).UTC(),
	}

	// When
	err := suite.Store.Set(ctx, shortURL1)
	require.NoError(t, err)
	err = suite.Store.Set(ctx, shortURL2)
	require.NoError(t, err)

	// Then
	retrievedURL, err := suite.Store.Get(ctx, shortURL1.Slug)
	require.NoError(t, err)
	assert.Equal(t, shortURL1.Slug, retrievedURL.Slug)
	assert.Equal(t, shortURL1.OriginalURL, retrievedURL.OriginalURL)
	if _, ok := suite.Store.(*CacheStore); !ok { // This assertion shouldn't be tested for cache
		assert.Equal(t, shortURL2.InsertedAt.Truncate(time.Second), retrievedURL.InsertedAt.Truncate(time.Second))
	}
}

func (suite *StoreTestSuite) TestDeleteExpired(t *testing.T) {
	// Given
	ctx := context.Background()
	shortURLExpired := domain.URLMapping{
		Slug:        "expired",
		OriginalURL: "https://example.com/expired",
		InsertedAt:  time.Now().UTC().Add(-24 * time.Hour),
	}
	shortURL := domain.URLMapping{
		Slug:        "active",
		OriginalURL: "https://example.com/active",
		InsertedAt:  time.Now().UTC(),
	}

	err := suite.Store.Set(ctx, shortURLExpired)
	require.NoError(t, err)
	err = suite.Store.Set(ctx, shortURL)
	require.NoError(t, err)

	// When
	slugsDeleted, err := suite.Store.DeleteExpired(ctx, 10*time.Hour)
	require.NoError(t, err)

	// Then
	require.Len(t, slugsDeleted, 1)
	assert.Equal(t, shortURLExpired.Slug, slugsDeleted[0])
	_, err = suite.Store.Get(ctx, shortURLExpired.Slug)
	assert.ErrorIs(t, err, ErrNotFound)
	_, err = suite.Store.Get(ctx, shortURL.Slug)
	assert.NoError(t, err)
}
