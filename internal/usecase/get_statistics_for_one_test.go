package usecase

import (
	"context"
	"testing"
	"urlShortenerService/domain"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/infrastructure/statistics"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetStatisticsForOneCmdBuilder(t *testing.T) {
	urlSanitizerStub := func(expectedURL *string, returnedURL string, err error) command.URLSanitizerCmd {
		return func(rawURL string) (string, error) {
			if expectedURL != nil {
				assert.Equal(t, *expectedURL, rawURL)
			}
			return returnedURL, err
		}
	}
	var originalURL string = "https://My-Very-Long-URL.com/needs-to-be-shortened"
	var sanitizedURL string = "https://my-very-long-url.com/needs-to-be-shortened"

	t.Run("nominal", func(t *testing.T) {
		// Given
		expectedURLStatistics := domain.URLStatistic{URL: sanitizedURL, AccessedCounter: 1, ShortenedCounter: 2}
		urlSanitizerCmd := urlSanitizerStub(&originalURL, sanitizedURL, nil)
		statisticsMock := statistics.NewMockStore(t)
		statisticsMock.On("GetOne", mock.Anything, sanitizedURL).Return(expectedURLStatistics, nil)
		cmd := GetStatisticsForOneCmdBuilder(urlSanitizerCmd, statisticsMock)

		// When
		urlStatisticsResp, err := cmd(context.Background(), originalURL)
		require.NoError(t, err)

		// Then
		assert.Equal(t, expectedURLStatistics, urlStatisticsResp)
	})
	t.Run("failed sanitizing URL", func(t *testing.T) {
		// Given
		urlSanitizerCmd := urlSanitizerStub(nil, "", assert.AnError)
		statisticsMock := statistics.NewMockStore(t)
		cmd := GetStatisticsForOneCmdBuilder(urlSanitizerCmd, statisticsMock)

		// When
		urlStatisticsResp, err := cmd(context.Background(), originalURL)

		// Then
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, urlStatisticsResp)
	})
	t.Run("failed updating statistics", func(t *testing.T) {
		// Given
		urlSanitizerCmd := urlSanitizerStub(&originalURL, sanitizedURL, nil)
		statisticsMock := statistics.NewMockStore(t)
		statisticsMock.On("GetOne", mock.Anything, mock.Anything).Return(domain.URLStatistic{}, assert.AnError)
		cmd := GetStatisticsForOneCmdBuilder(urlSanitizerCmd, statisticsMock)

		// When
		urlStatisticsResp, err := cmd(context.Background(), originalURL)

		// Then
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, urlStatisticsResp)
	})
}
