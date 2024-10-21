package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"urlShortenerService/domain"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/infrastructure/shorturl"
	"urlShortenerService/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithGetOriginalURLForceHandler(t *testing.T) {
	originalURL := "https://my-very-long-url.com/needs-to-be-shortened"
	slug := "zTw34enA"
	mockCmd := func(err error) usecase.GetOriginalURLCmd {
		return func(ctx context.Context, s string) (string, error) {
			assert.Equal(t, slug, s)
			return originalURL, err
		}
	}

	t.Run("ok", func(t *testing.T) {
		// Given
		router := NewBuilder(domain.EnvTest).WithGetOriginalURLForceHandler(mockCmd(nil)).router
		u, err := url.Parse(fmt.Sprintf("/%s/force", slug))
		require.NoError(t, err)

		// When
		record := httptest.NewRecorder()
		req := httptest.NewRequest("GET", u.String(), nil)
		router.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusOK, record.Code)
		bodyResponse := GetOriginalURLResponse{}
		require.NoError(t, json.Unmarshal(record.Body.Bytes(), &bodyResponse))
		assert.Equal(t, originalURL, bodyResponse.OriginalURL)
	})
	t.Run("redirection asked", func(t *testing.T) {
		// Given
		router := NewBuilder(domain.EnvTest).WithGetOriginalURLForceHandler(mockCmd(nil)).router
		u, err := url.Parse(fmt.Sprintf("/%s/force?redirect=true", slug))
		require.NoError(t, err)

		// When
		record := httptest.NewRecorder()
		req := httptest.NewRequest("GET", u.String(), nil)
		router.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusFound, record.Code)
		loc, err := record.Result().Location()
		require.NoError(t, err)
		assert.Equal(t, originalURL, loc.String())
	})
	t.Run("redirection false", func(t *testing.T) {
		// Given
		router := NewBuilder(domain.EnvTest).WithGetOriginalURLForceHandler(mockCmd(nil)).router
		u, err := url.Parse(fmt.Sprintf("/%s/force?redirect=false", slug))
		require.NoError(t, err)

		// When
		record := httptest.NewRecorder()
		req := httptest.NewRequest("GET", u.String(), nil)
		router.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusOK, record.Code)
		bodyResponse := GetOriginalURLResponse{}
		require.NoError(t, json.Unmarshal(record.Body.Bytes(), &bodyResponse))
		assert.Equal(t, originalURL, bodyResponse.OriginalURL)
	})
	t.Run("redirection invalid", func(t *testing.T) {
		// Given
		router := NewBuilder(domain.EnvTest).WithGetOriginalURLForceHandler(mockCmd(nil)).router
		u, err := url.Parse(fmt.Sprintf("/%s/force?redirect=something", slug))
		require.NoError(t, err)

		// When
		record := httptest.NewRecorder()
		req := httptest.NewRequest("GET", u.String(), nil)
		router.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusOK, record.Code)
		bodyResponse := GetOriginalURLResponse{}
		require.NoError(t, json.Unmarshal(record.Body.Bytes(), &bodyResponse))
		assert.Equal(t, originalURL, bodyResponse.OriginalURL)
	})
	t.Run("not found", func(t *testing.T) {
		// Given
		router := NewBuilder(domain.EnvTest).WithGetOriginalURLForceHandler(mockCmd(shorturl.ErrNotFound)).router
		u, err := url.Parse(fmt.Sprintf("/%s/force?redirect=true", slug))
		require.NoError(t, err)

		// When
		record := httptest.NewRecorder()
		req := httptest.NewRequest("GET", u.String(), nil)
		router.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusNotFound, record.Code)
	})
	t.Run("unprocessable entity", func(t *testing.T) {
		t.Run("invalid slug lenght", func(t *testing.T) {
			// Given
			router := NewBuilder(domain.EnvTest).WithGetOriginalURLForceHandler(mockCmd(command.ErrInvalidSlugLenght)).router
			u, err := url.Parse(fmt.Sprintf("/%s/force", slug))
			require.NoError(t, err)

			// When
			record := httptest.NewRecorder()
			req := httptest.NewRequest("GET", u.String(), nil)
			router.ServeHTTP(record, req)

			// Then
			assert.Equal(t, http.StatusUnprocessableEntity, record.Code)
		})
		t.Run("invalid slug with non alphanumeric character", func(t *testing.T) {
			// Given
			router := NewBuilder(domain.EnvTest).WithGetOriginalURLForceHandler(mockCmd(command.ErrInvalidSlugNonAlphanumeric)).router
			u, err := url.Parse(fmt.Sprintf("/%s/force", slug))
			require.NoError(t, err)

			// When
			record := httptest.NewRecorder()
			req := httptest.NewRequest("GET", u.String(), nil)
			router.ServeHTTP(record, req)

			// Then
			assert.Equal(t, http.StatusUnprocessableEntity, record.Code)
		})
	})
	t.Run("internal server error", func(t *testing.T) {
		// Given
		router := NewBuilder(domain.EnvTest).WithGetOriginalURLForceHandler(mockCmd(assert.AnError)).router
		u, err := url.Parse(fmt.Sprintf("/%s/force", slug))
		require.NoError(t, err)

		// When
		record := httptest.NewRecorder()
		req := httptest.NewRequest("GET", u.String(), nil)
		router.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusInternalServerError, record.Code)
	})
}
