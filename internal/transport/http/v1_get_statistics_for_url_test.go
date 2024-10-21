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
	"urlShortenerService/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithGetStatisticsForURLHandler(t *testing.T) {
	urlToStat := "https://example.com"
	urlStatistics := domain.URLStatistic{
		URL:              urlToStat,
		ShortenedCounter: 1,
		AccessedCounter:  10,
	}
	mockCmd := func(expectedURL *string, urlStatistics domain.URLStatistic, err error) usecase.GetStatisticsForURLCmd {
		return func(ctx context.Context, urlToStat string) (domain.URLStatistic, error) {
			if expectedURL != nil {
				assert.Equal(t, *expectedURL, urlToStat)
			}
			return urlStatistics, err
		}
	}

	t.Run("ok", func(t *testing.T) {
		// Given
		router := NewBuilder(domain.EnvTest).WithGetStatisticsForURLHandler(mockCmd(&urlToStat, urlStatistics, nil)).router
		u, err := url.Parse(fmt.Sprintf("%s/statistics?encoded_url=%s", pathPrefixV1, url.QueryEscape(urlToStat)))
		require.NoError(t, err)

		// When
		record := httptest.NewRecorder()
		req := httptest.NewRequest("GET", u.String(), nil)
		router.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusOK, record.Code)
		bodyResponse := GetStatisticsForURLResponse{}
		require.NoError(t, json.Unmarshal(record.Body.Bytes(), &bodyResponse))
		assert.Equal(t, urlStatistics.URL, bodyResponse.URL)
		assert.Equal(t, urlStatistics.AccessedCounter, bodyResponse.AccessedCounter)
		assert.Equal(t, urlStatistics.ShortenedCounter, bodyResponse.ShortenedCounter)
	})
	t.Run("bad request", func(t *testing.T) {
		t.Run("missing encoded_url query parameter", func(t *testing.T) {
			// Given
			router := NewBuilder(domain.EnvTest).WithGetStatisticsForURLHandler(mockCmd(nil, urlStatistics, nil)).router
			u, err := url.Parse(fmt.Sprintf("%s/statistics", pathPrefixV1))
			require.NoError(t, err)

			// When
			record := httptest.NewRecorder()
			req := httptest.NewRequest("GET", u.String(), nil)
			router.ServeHTTP(record, req)

			// Then
			assert.Equal(t, http.StatusBadRequest, record.Code)
		})
		t.Run("badly encoded URL", func(t *testing.T) {
			// Given
			router := NewBuilder(domain.EnvTest).WithGetStatisticsForURLHandler(mockCmd(&urlToStat, urlStatistics, nil)).router
			u, err := url.Parse(fmt.Sprintf("%s/statistics?encoded_url=https://badly-encoded.com%%", pathPrefixV1))
			require.NoError(t, err)

			// When
			record := httptest.NewRecorder()
			req := httptest.NewRequest("GET", u.String(), nil)
			router.ServeHTTP(record, req)

			// Then
			assert.Equal(t, http.StatusBadRequest, record.Code)
		})
	})
	t.Run("internal server error", func(t *testing.T) {
		// Given
		router := NewBuilder(domain.EnvTest).WithGetStatisticsForURLHandler(mockCmd(&urlToStat, urlStatistics, assert.AnError)).router
		u, err := url.Parse(fmt.Sprintf("%s/statistics?encoded_url=%s", pathPrefixV1, url.QueryEscape(urlToStat)))
		require.NoError(t, err)

		// When
		record := httptest.NewRecorder()
		req := httptest.NewRequest("GET", u.String(), nil)
		router.ServeHTTP(record, req)

		// Then
		assert.Equal(t, http.StatusInternalServerError, record.Code)
	})
}
