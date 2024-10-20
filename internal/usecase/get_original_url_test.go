package usecase

import (
	"context"
	"testing"
	"urlShortenerService/domain"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/infrastructure/shorturl"

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
		cmd := GetOriginalURLCmdBuilder(slugValidatorCmd, shortURLMock)

		// When
		originalURL, err := cmd(context.Background(), urlMappingData.Slug)
		require.NoError(t, err)

		// Then
		assert.Equal(t, urlMappingData.OriginalURL, originalURL)
	})
	t.Run("invalid slug", func(t *testing.T) {
		// Given
		slugValidatorCmd := slugValidatorStub(nil, assert.AnError)
		shortURLMock := shorturl.NewMock(t)
		cmd := GetOriginalURLCmdBuilder(slugValidatorCmd, shortURLMock)

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
		cmd := GetOriginalURLCmdBuilder(slugValidatorCmd, shortURLMock)

		// When
		originalURL, err := cmd(context.Background(), urlMappingData.Slug)

		// Then
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, originalURL)
	})
}
