package usecase

import (
	"context"
	"sync"
	"testing"
	"urlShortenerService/domain"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/infrastructure/shorturl"
	"urlShortenerService/internal/infrastructure/statistics"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetOriginalURLCmdBuilder(t *testing.T) {
	slugValidatorStub := func(expectedSlug *string, err error) command.SlugValidatorCmd {
		return func(slug string) error {
			if expectedSlug != nil {
				assert.Equal(t, *expectedSlug, slug)
			}
			return err
		}
	}
	var urlMappingData domain.URLMapping = domain.URLMapping{
		Slug:        "zTw34enA",
		OriginalURL: "https://My-Very-Long-URL.com/needs-to-be-shortened",
	}

	t.Run("nominal", func(t *testing.T) {
		// Given
		slugValidatorCmd := slugValidatorStub(&urlMappingData.Slug, nil)
		shortURLMock := shorturl.NewMock(t)
		shortURLMock.On("Get", mock.Anything, urlMappingData.Slug).Return(urlMappingData, nil)
		var wg sync.WaitGroup
		wg.Add(1)
		statisticsMock := statistics.NewMockStore(t)
		statisticsMock.On("Set", mock.Anything, urlMappingData.OriginalURL, statistics.StatisticTypeAccessed).Return(nil).Run(func(args mock.Arguments) {
			wg.Done()
		})
		cmd := GetOriginalURLCmdBuilder(slugValidatorCmd, shortURLMock, statisticsMock)

		// When
		originalURL, err := cmd(context.Background(), urlMappingData.Slug)
		require.NoError(t, err)

		// Then
		assert.Equal(t, urlMappingData.OriginalURL, originalURL)
		wg.Wait()
	})
	t.Run("invalid slug", func(t *testing.T) {
		// Given
		slugValidatorCmd := slugValidatorStub(nil, assert.AnError)
		shortURLMock := shorturl.NewMock(t)
		statisticsMock := statistics.NewMockStore(t)
		cmd := GetOriginalURLCmdBuilder(slugValidatorCmd, shortURLMock, statisticsMock)

		// When
		originalURL, err := cmd(context.Background(), urlMappingData.Slug)

		// Then
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, originalURL)
	})
	t.Run("failed retrieving URL", func(t *testing.T) {
		// Given
		slugValidatorCmd := slugValidatorStub(nil, nil)
		shortURLMock := shorturl.NewMock(t)
		shortURLMock.On("Get", mock.Anything, mock.Anything).Return(domain.URLMapping{}, assert.AnError)
		statisticsMock := statistics.NewMockStore(t)
		cmd := GetOriginalURLCmdBuilder(slugValidatorCmd, shortURLMock, statisticsMock)

		// When
		originalURL, err := cmd(context.Background(), urlMappingData.Slug)

		// Then
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, originalURL)
	})
	t.Run("failed updating statistics", func(t *testing.T) {
		// Given
		slugValidatorCmd := slugValidatorStub(&urlMappingData.Slug, nil)
		shortURLMock := shorturl.NewMock(t)
		shortURLMock.On("Get", mock.Anything, urlMappingData.Slug).Return(urlMappingData, nil)
		var wg sync.WaitGroup
		wg.Add(1)
		statisticsMock := statistics.NewMockStore(t)
		statisticsMock.On("Set", mock.Anything, urlMappingData.OriginalURL, statistics.StatisticTypeAccessed).Return(assert.AnError).Run(func(args mock.Arguments) {
			wg.Done()
		})
		cmd := GetOriginalURLCmdBuilder(slugValidatorCmd, shortURLMock, statisticsMock)

		// When
		originalURL, err := cmd(context.Background(), urlMappingData.Slug)
		require.NoError(t, err)

		// Then
		assert.Equal(t, urlMappingData.OriginalURL, originalURL)
		wg.Wait()
	})
}
