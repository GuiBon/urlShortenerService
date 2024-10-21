package usecase

import (
	"context"
	"testing"
	"urlShortenerService/domain"
	"urlShortenerService/internal/infrastructure/statistics"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetTopStatisticsCmdBuilder(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		// Given
		var statType statistics.StatisticType = statistics.StatisticTypeShortened
		var limitOveride int64 = 3
		expectedURLStatistics := []domain.URLStatistic{
			{URL: "https://example.com/1", ShortenedCounter: 2},
			{URL: "https://example.com/2", ShortenedCounter: 4},
		}
		statisticsMock := statistics.NewMockStore(t)
		statisticsMock.On("GetTopURLs", mock.Anything, statType, limitOveride).Return(expectedURLStatistics, nil)
		cmd := GetTopStatisticsCmdBuilder(statisticsMock)

		// When
		urlStatisticsResp, err := cmd(context.Background(), statType, limitOveride)
		require.NoError(t, err)

		// Then
		assert.Equal(t, expectedURLStatistics, urlStatisticsResp)
	})
	t.Run("failed retrieving statistics", func(t *testing.T) {
		// Given
		statisticsMock := statistics.NewMockStore(t)
		statisticsMock.On("GetTopURLs", mock.Anything, mock.Anything, mock.Anything).Return([]domain.URLStatistic{}, assert.AnError)
		cmd := GetTopStatisticsCmdBuilder(statisticsMock)

		// When
		urlStatisticsResp, err := cmd(context.Background(), statistics.StatisticTypeShortened, 0)

		// Then
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, urlStatisticsResp)
	})
}
