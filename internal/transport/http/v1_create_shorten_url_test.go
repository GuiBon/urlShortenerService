package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"urlShortenerService/domain"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithV1CreateShortenURLHandler(t *testing.T) {
	originalURL := "https://my-very-long-url.com/needs-to-be-shortened"
	shortURL := "https://localhost:8080/zTw34enA"
	u, err := url.Parse(fmt.Sprintf("%s/shorten", pathPrefixV1))
	require.NoError(t, err)
	mockCmd := func(err error) usecase.CreateShortenURLCmd {
		return func(ctx context.Context, urlToShorten string) (string, error) {
			assert.Equal(t, originalURL, urlToShorten)
			return shortURL, err
		}
	}

	t.Run("created", func(t *testing.T) {
		// Given
		router := NewBuilder(domain.EnvTest).WithV1CreateShortenURLHandler(mockCmd(nil)).router

		// When
		record := httptest.NewRecorder()
		req := httptest.NewRequest("POST", u.String(), strings.NewReader(fmt.Sprintf(`{"original_url": "%s"}`, originalURL)))
		router.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusCreated, record.Code)
		bodyResponse := CreateShortenURLResponse{}
		require.NoError(t, json.Unmarshal(record.Body.Bytes(), &bodyResponse))
		assert.Equal(t, shortURL, bodyResponse.ShortURL)
	})
	t.Run("bad request", func(t *testing.T) {
		t.Run("not a valid JSON body", func(t *testing.T) {
			// Given
			router := NewBuilder(domain.EnvTest).WithV1CreateShortenURLHandler(mockCmd(nil)).router

			// When
			record := httptest.NewRecorder()
			req := httptest.NewRequest("POST", u.String(), strings.NewReader(`{`))
			router.ServeHTTP(record, req)

			// Then
			assert.Equal(t, http.StatusBadRequest, record.Code)
		})
		t.Run("missing required field in JSON body", func(t *testing.T) {
			// Given
			router := NewBuilder(domain.EnvTest).WithV1CreateShortenURLHandler(mockCmd(nil)).router

			// When
			record := httptest.NewRecorder()
			req := httptest.NewRequest("POST", u.String(), strings.NewReader(`{}`))
			router.ServeHTTP(record, req)

			// Then
			assert.Equal(t, http.StatusBadRequest, record.Code)
		})
	})
	t.Run("unprocessable entity", func(t *testing.T) {
		// Given
		router := NewBuilder(domain.EnvTest).WithV1CreateShortenURLHandler(mockCmd(command.ErrInvalidURL)).router

		// When
		record := httptest.NewRecorder()
		req := httptest.NewRequest("POST", u.String(), strings.NewReader(fmt.Sprintf(`{"original_url": "%s"}`, originalURL)))
		router.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusUnprocessableEntity, record.Code)
	})
	t.Run("internal server error", func(t *testing.T) {
		// Given
		router := NewBuilder(domain.EnvTest).WithV1CreateShortenURLHandler(mockCmd(assert.AnError)).router

		// When
		record := httptest.NewRecorder()
		req := httptest.NewRequest("POST", u.String(), strings.NewReader(fmt.Sprintf(`{"original_url": "%s"}`, originalURL)))
		router.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusInternalServerError, record.Code)
	})
}
