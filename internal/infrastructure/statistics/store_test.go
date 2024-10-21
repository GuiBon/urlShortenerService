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

	t.Run("TestSet", suite.TestSet)
	t.Run("TestGetOne", suite.TestGetOne)
	t.Run("TestGetTop", suite.TestGetTop)
}

func (suite *StoreTestSuite) TestSet(t *testing.T) {
	// Given
	ctx := context.Background()
	url := "https://example.com/set-test"

	// When
	err := suite.Store.Set(ctx, url, StatisticTypeAccessed)
	require.NoError(t, err)
	err = suite.Store.Set(ctx, url, StatisticTypeAccessed)
	require.NoError(t, err)
	err = suite.Store.Set(ctx, url, StatisticTypeShortened)
	require.NoError(t, err)

	// Then
	stats, err := suite.Store.GetOne(ctx, url)
	require.NoError(t, err)
	assert.Equal(t, url, stats.URL)
	assert.Equal(t, 2, stats.AccessedCounter)
	assert.Equal(t, 1, stats.ShortenedCounter)
}

func (suite *StoreTestSuite) TestGetOne(t *testing.T) {
	// Given
	ctx := context.Background()
	url := "https://example.com/getone-test"
	err := suite.Store.Set(ctx, url, StatisticTypeAccessed)
	require.NoError(t, err)
	err = suite.Store.Set(ctx, url, StatisticTypeShortened)
	require.NoError(t, err)

	// When
	stats, err := suite.Store.GetOne(ctx, url)
	require.NoError(t, err)

	// Then
	assert.Equal(t, url, stats.URL)
	assert.Equal(t, 1, stats.AccessedCounter)
	assert.Equal(t, 1, stats.ShortenedCounter)
}

func (suite *StoreTestSuite) TestGetTop(t *testing.T) {
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
				err := suite.Store.Set(ctx, stat.URL, StatisticTypeAccessed)
				require.NoError(t, err)
			}
		}

		// When
		stats, err := suite.Store.GetTop(ctx, StatisticTypeAccessed, int64(len(expectedStats)))
		require.NoError(t, err)

		// Then
		assert.Equal(t, expectedStats, stats)
	})
}
