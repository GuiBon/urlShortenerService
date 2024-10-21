package statistics

import (
	context "context"
	"testing"
	"urlShortenerService/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StoreTestSuite struct {
	Store
}

func RunStoreTests(t *testing.T, store Store) {
	suite := &StoreTestSuite{Store: store}

	t.Run("TestSetURL", suite.TestSetURL)
	t.Run("TestGetURL", suite.TestGetURL)
	t.Run("TestGetTopURLs", suite.TestGetTopURLs)
}

func (suite *StoreTestSuite) TestSetURL(t *testing.T) {
	// Given
	ctx := context.Background()
	url := "https://example.com/set-test"

	// When
	err := suite.Store.SetURL(ctx, url, StatisticTypeAccessed)
	require.NoError(t, err)
	err = suite.Store.SetURL(ctx, url, StatisticTypeAccessed)
	require.NoError(t, err)
	err = suite.Store.SetURL(ctx, url, StatisticTypeShortened)
	require.NoError(t, err)

	// Then
	stats, err := suite.Store.GetURL(ctx, url)
	require.NoError(t, err)
	assert.Equal(t, url, stats.URL)
	assert.Equal(t, 2, stats.AccessedCounter)
	assert.Equal(t, 1, stats.ShortenedCounter)
}

func (suite *StoreTestSuite) TestGetURL(t *testing.T) {
	// Given
	ctx := context.Background()
	url := "https://example.com/getone-test"
	err := suite.Store.SetURL(ctx, url, StatisticTypeAccessed)
	require.NoError(t, err)
	err = suite.Store.SetURL(ctx, url, StatisticTypeShortened)
	require.NoError(t, err)

	// When
	stats, err := suite.Store.GetURL(ctx, url)
	require.NoError(t, err)

	// Then
	assert.Equal(t, url, stats.URL)
	assert.Equal(t, 1, stats.AccessedCounter)
	assert.Equal(t, 1, stats.ShortenedCounter)
}

func (suite *StoreTestSuite) TestGetTopURLs(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		// Given
		ctx := context.Background()
		expectedStats := []domain.URLStatistic{
			{URL: "https://example.com/gettop-test-3", AccessedCounter: 15},
			{URL: "https://example.com/gettop-test-1", AccessedCounter: 10},
			{URL: "https://example.com/gettop-test-2", AccessedCounter: 5},
		}
		for _, stat := range expectedStats {
			for i := 0; i < stat.AccessedCounter; i++ {
				err := suite.Store.SetURL(ctx, stat.URL, StatisticTypeAccessed)
				require.NoError(t, err)
			}
		}

		// When
		stats, err := suite.Store.GetTopURLs(ctx, StatisticTypeAccessed, int64(len(expectedStats)))
		require.NoError(t, err)

		// Then
		assert.Equal(t, expectedStats, stats)
	})
}
