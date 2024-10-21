package usecase

import (
	"context"
	"sync"
	"testing"
	"urlShortenerService/domain"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/infrastructure/malwarescanner"
	"urlShortenerService/internal/infrastructure/shorturl"
	"urlShortenerService/internal/infrastructure/statistics"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetOriginalURLWithMalwareScanCmdBuilder(t *testing.T) {
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
		OriginalURL: "https://My-Very-Long-URL.com/needs-to-be-shortened/malware",
	}

	t.Run("nominal", func(t *testing.T) {
		// Given
		slugValidatorCmd := slugValidatorStub(&urlMappingData.Slug, nil)
		malwareScannerMock := malwarescanner.NewScannerMock(t)
		malwareScannerMock.On("Scan", mock.Anything, urlMappingData.OriginalURL, mock.Anything).Return(malwarescanner.MalwareScanResultClear)
		shortURLMock := shorturl.NewMock(t)
		shortURLMock.On("Get", mock.Anything, urlMappingData.Slug).Return(urlMappingData, nil)
		var wg sync.WaitGroup
		wg.Add(1)
		statisticsMock := statistics.NewMockStore(t)
		statisticsMock.On("SetURL", mock.Anything, urlMappingData.OriginalURL, statistics.StatisticTypeAccessed).Return(nil).Run(func(args mock.Arguments) {
			wg.Done()
		})
		cmd := GetOriginalURLWithMalwareScanCmdBuilder(slugValidatorCmd, malwareScannerMock, shortURLMock, statisticsMock)

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
		malwareScannerMock := malwarescanner.NewScannerMock(t)
		shortURLMock := shorturl.NewMock(t)
		statisticsMock := statistics.NewMockStore(t)
		cmd := GetOriginalURLWithMalwareScanCmdBuilder(slugValidatorCmd, malwareScannerMock, shortURLMock, statisticsMock)

		// When
		originalURL, err := cmd(context.Background(), urlMappingData.Slug)

		// Then
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, originalURL)
	})
	t.Run("failed retrieving URL", func(t *testing.T) {
		// Given
		slugValidatorCmd := slugValidatorStub(nil, nil)
		malwareScannerMock := malwarescanner.NewScannerMock(t)
		shortURLMock := shorturl.NewMock(t)
		shortURLMock.On("Get", mock.Anything, mock.Anything).Return(domain.URLMapping{}, assert.AnError)
		statisticsMock := statistics.NewMockStore(t)
		cmd := GetOriginalURLWithMalwareScanCmdBuilder(slugValidatorCmd, malwareScannerMock, shortURLMock, statisticsMock)

		// When
		originalURL, err := cmd(context.Background(), urlMappingData.Slug)

		// Then
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, originalURL)
	})
	t.Run("failed updating statistics", func(t *testing.T) {
		// Given
		slugValidatorCmd := slugValidatorStub(&urlMappingData.Slug, nil)
		malwareScannerMock := malwarescanner.NewScannerMock(t)
		malwareScannerMock.On("Scan", mock.Anything, urlMappingData.OriginalURL, mock.Anything).Return(malwarescanner.MalwareScanResultClear)
		shortURLMock := shorturl.NewMock(t)
		shortURLMock.On("Get", mock.Anything, urlMappingData.Slug).Return(urlMappingData, nil)
		var wg sync.WaitGroup
		wg.Add(1)
		statisticsMock := statistics.NewMockStore(t)
		statisticsMock.On("SetURL", mock.Anything, urlMappingData.OriginalURL, statistics.StatisticTypeAccessed).Return(assert.AnError).Run(func(args mock.Arguments) {
			wg.Done()
		})
		cmd := GetOriginalURLWithMalwareScanCmdBuilder(slugValidatorCmd, malwareScannerMock, shortURLMock, statisticsMock)

		// When
		originalURL, err := cmd(context.Background(), urlMappingData.Slug)
		require.NoError(t, err)

		// Then
		assert.Equal(t, urlMappingData.OriginalURL, originalURL)
		wg.Wait()
	})
	t.Run("failed to scan the URL for malware", func(t *testing.T) {
		// Given
		slugValidatorCmd := slugValidatorStub(&urlMappingData.Slug, nil)
		malwareScannerMock := malwarescanner.NewScannerMock(t)
		malwareScannerMock.On("Scan", mock.Anything, urlMappingData.OriginalURL, mock.Anything).Return(malwarescanner.MalwareScanUnknownError)
		shortURLMock := shorturl.NewMock(t)
		shortURLMock.On("Get", mock.Anything, urlMappingData.Slug).Return(urlMappingData, nil)
		var wg sync.WaitGroup
		wg.Add(1)
		statisticsMock := statistics.NewMockStore(t)
		statisticsMock.On("SetURL", mock.Anything, urlMappingData.OriginalURL, statistics.StatisticTypeAccessed).Return(nil).Run(func(args mock.Arguments) {
			wg.Done()
		})
		cmd := GetOriginalURLWithMalwareScanCmdBuilder(slugValidatorCmd, malwareScannerMock, shortURLMock, statisticsMock)

		// When
		originalURL, err := cmd(context.Background(), urlMappingData.Slug)
		require.NoError(t, err)

		// Then
		assert.Equal(t, urlMappingData.OriginalURL, originalURL)
		wg.Wait()
	})
	t.Run("malware detected", func(t *testing.T) {
		// Given
		slugValidatorCmd := slugValidatorStub(&urlMappingData.Slug, nil)
		malwareScannerMock := malwarescanner.NewScannerMock(t)
		malwareScannerMock.On("Scan", mock.Anything, urlMappingData.OriginalURL, mock.Anything).Return(malwarescanner.MalwareScanResultDetected)
		shortURLMock := shorturl.NewMock(t)
		shortURLMock.On("Get", mock.Anything, urlMappingData.Slug).Return(urlMappingData, nil)
		var wg sync.WaitGroup
		wg.Add(1)
		statisticsMock := statistics.NewMockStore(t)
		statisticsMock.On("SetURL", mock.Anything, urlMappingData.OriginalURL, statistics.StatisticTypeAccessed).Return(nil).Run(func(args mock.Arguments) {
			wg.Done()
		})
		cmd := GetOriginalURLWithMalwareScanCmdBuilder(slugValidatorCmd, malwareScannerMock, shortURLMock, statisticsMock)

		// When
		originalURL, err := cmd(context.Background(), urlMappingData.Slug)

		// Then
		require.ErrorIs(t, err, malwarescanner.ErrMalswareURL)
		assert.Empty(t, originalURL)
		wg.Wait()
	})
}

func TestForceGetOriginalURLCmdBuilder(t *testing.T) {
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
		statisticsMock.On("SetURL", mock.Anything, urlMappingData.OriginalURL, statistics.StatisticTypeAccessed).Return(nil).Run(func(args mock.Arguments) {
			wg.Done()
		})
		cmd := ForceGetOriginalURLCmdBuilder(slugValidatorCmd, shortURLMock, statisticsMock)

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
		cmd := ForceGetOriginalURLCmdBuilder(slugValidatorCmd, shortURLMock, statisticsMock)

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
		cmd := ForceGetOriginalURLCmdBuilder(slugValidatorCmd, shortURLMock, statisticsMock)

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
		statisticsMock.On("SetURL", mock.Anything, urlMappingData.OriginalURL, statistics.StatisticTypeAccessed).Return(assert.AnError).Run(func(args mock.Arguments) {
			wg.Done()
		})
		cmd := ForceGetOriginalURLCmdBuilder(slugValidatorCmd, shortURLMock, statisticsMock)

		// When
		originalURL, err := cmd(context.Background(), urlMappingData.Slug)
		require.NoError(t, err)

		// Then
		assert.Equal(t, urlMappingData.OriginalURL, originalURL)
		wg.Wait()
	})
}
