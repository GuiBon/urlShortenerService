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
	shortURL := domain.ShortURL{
		Slug: "example",
		URL:  "https://example.com",
	}

	// When
	err := suite.Store.Set(ctx, shortURL)
	require.NoError(t, err)

	// Then
	retrievedURL, err := suite.Store.Get(ctx, shortURL.Slug)
	require.NoError(t, err)
	assert.Equal(t, shortURL.Slug, retrievedURL.Slug)
	assert.Equal(t, shortURL.URL, retrievedURL.URL)
	assert.NotEmpty(t, retrievedURL.InsertedAt)
}

func (suite *StoreTestSuite) TestGet(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		// Given
		ctx := context.Background()
		shortURL := domain.ShortURL{
			Slug: "example",
			URL:  "https://example.com",
		}
		err := suite.Store.Set(ctx, shortURL)
		require.NoError(t, err)

		// When
		retrievedURL, err := suite.Store.Get(ctx, shortURL.Slug)
		require.NoError(t, err)

		// Then
		assert.Equal(t, shortURL.Slug, retrievedURL.Slug)
		assert.Equal(t, shortURL.URL, retrievedURL.URL)
		assert.NotEmpty(t, retrievedURL.InsertedAt)
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
	shortURL1 := domain.ShortURL{
		Slug:       "duplicate",
		URL:        "https://example.com/duplicate",
		InsertedAt: time.Now().Add(-4 * time.Hour).UTC(),
	}
	shortURL2 := domain.ShortURL{
		Slug:       "duplicate",
		URL:        "https://example.com/duplicate",
		InsertedAt: time.Now().Add(-1 * time.Hour).UTC(),
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
	assert.Equal(t, shortURL1.URL, retrievedURL.URL)
	assert.Equal(t, shortURL2.InsertedAt.Truncate(time.Second), retrievedURL.InsertedAt.Truncate(time.Second))
}

func (suite *StoreTestSuite) TestDeleteExpired(t *testing.T) {
	// Given
	ctx := context.Background()
	shortURL := domain.ShortURL{
		Slug:       "expired",
		URL:        "https://example.com/expired",
		InsertedAt: time.Now().Add(-24 * time.Hour).UTC(),
	}
	err := suite.Store.Set(ctx, shortURL)
	require.NoError(t, err)

	// When
	err = suite.Store.DeleteExpired(ctx, 1*time.Hour)
	require.NoError(t, err)

	// Then
	_, err = suite.Store.Get(ctx, shortURL.Slug)
	assert.ErrorIs(t, err, ErrNotFound)
}
