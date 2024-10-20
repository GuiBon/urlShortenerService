package usecase

import (
	"context"
	"fmt"
	"testing"
	"urlShortenerService/domain"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/infrastructure/shorturl"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateShortenURLCmdBuilder(t *testing.T) {
	urlSanitizerStub := func(expectedURL *string, returnedURL string, err error) command.URLSanitizerCmd {
		return func(rawURL string) (string, error) {
			if expectedURL != nil {
				assert.Equal(t, *expectedURL, rawURL)
			}
			return returnedURL, err
		}
	}
	slugGeneratorStub := func(expectedURL *string, returnedSlug string) command.SlugGeneratorCmd {
		return func(rawURL string) string {
			if expectedURL != nil {
				assert.Equal(t, *expectedURL, rawURL)
			}
			return returnedSlug
		}
	}
	var baseURL string = "https://example.com"
	var originalURL string = "https://My-Very-Long-URL.com/needs-to-be-shortened"
	var sanitizedURL string = "https://my-very-long-url.com/needs-to-be-shortened"
	var slug string = "zTw34enA"

	t.Run("nominal", func(t *testing.T) {
		// Given
		urlSanitizerCmd := urlSanitizerStub(&originalURL, sanitizedURL, nil)
		slugGeneratorCmd := slugGeneratorStub(&sanitizedURL, slug)
		shortURLMock := shorturl.NewMock(t)
		shortURLMock.On("Set", mock.Anything, domain.ShortURL{Slug: slug, URL: sanitizedURL}).Return(nil)
		cmd := CreateShortenURLCmdBuilder(baseURL, urlSanitizerCmd, slugGeneratorCmd, shortURLMock)

		// When
		shortURL, err := cmd(context.Background(), originalURL)
		require.NoError(t, err)

		// Then
		assert.Equal(t, fmt.Sprintf("%s/%s", baseURL, slug), shortURL)
	})
	t.Run("failed sanitizing URL", func(t *testing.T) {
		// Given
		urlSanitizerCmd := urlSanitizerStub(nil, "", assert.AnError)
		slugGeneratorCmd := slugGeneratorStub(nil, slug)
		shortURLMock := shorturl.NewMock(t)
		cmd := CreateShortenURLCmdBuilder(baseURL, urlSanitizerCmd, slugGeneratorCmd, shortURLMock)

		// When
		shortURL, err := cmd(context.Background(), originalURL)

		// Then
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, shortURL)
	})
	t.Run("failed storing URL", func(t *testing.T) {
		// Given
		urlSanitizerCmd := urlSanitizerStub(nil, sanitizedURL, nil)
		slugGeneratorCmd := slugGeneratorStub(nil, slug)
		shortURLMock := shorturl.NewMock(t)
		shortURLMock.On("Set", mock.Anything, mock.Anything).Return(assert.AnError)
		cmd := CreateShortenURLCmdBuilder(baseURL, urlSanitizerCmd, slugGeneratorCmd, shortURLMock)

		// When
		shortURL, err := cmd(context.Background(), originalURL)

		// Then
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, shortURL)
	})
}
