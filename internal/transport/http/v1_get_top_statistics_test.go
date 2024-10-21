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
	"urlShortenerService/internal/infrastructure/statistics"
	"urlShortenerService/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithGetTopStatisticsHandler(t *testing.T) {
	var topAccessedStatistics []domain.URLStatistic = []domain.URLStatistic{
		{URL: "https://example.com/1", AccessedCounter: 10},
		{URL: "https://example.com/2", AccessedCounter: 5},
	}
	var topShortenedStatistics []domain.URLStatistic = []domain.URLStatistic{
		{URL: "https://example.com/1", ShortenedCounter: 4},
		{URL: "https://example.com/2", ShortenedCounter: 1},
	}
	var expectedTopStatisticsAccessedResponse GetTopStatisticsResponse = GetTopStatisticsResponse{
		URLs: []getTopStatisticsForURLResponse{
			{URL: "https://example.com/1", AccessedCounter: 10},
			{URL: "https://example.com/2", AccessedCounter: 5},
		},
	}
	var expectedTopStatisticsShortenedResponse GetTopStatisticsResponse = GetTopStatisticsResponse{
		URLs: []getTopStatisticsForURLResponse{
			{URL: "https://example.com/1", ShortenedCounter: 4},
			{URL: "https://example.com/2", ShortenedCounter: 1},
		},
	}
	var limit int64 = 10
	mockCmd := func(expectedStatType *statistics.StatisticType, expectedLimit int64, urlStatistics []domain.URLStatistic, err error) usecase.GetTopStatisticsCmd {
		return func(ctx context.Context, statType statistics.StatisticType, limit int64) ([]domain.URLStatistic, error) {
			if expectedStatType != nil {
				assert.Equal(t, *expectedStatType, statType)
			}
			assert.Equal(t, expectedLimit, limit)
			return urlStatistics, err
		}
	}

	t.Run("for accessed", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			// Given
			router := NewBuilder(domain.EnvTest).WithGetTopStatisticsHandler(mockCmd(&statistics.StatisticTypeAccessed, limit, topAccessedStatistics, nil)).router
			u, err := url.Parse(fmt.Sprintf("%s/statistics/accessed?limit=%d", pathPrefixV1, limit))
			require.NoError(t, err)

			// When
			record := httptest.NewRecorder()
			req := httptest.NewRequest("GET", u.String(), nil)
			router.ServeHTTP(record, req)

			// Then
			assert.Equal(t, http.StatusOK, record.Code)
			bodyResponse := GetTopStatisticsResponse{}
			require.NoError(t, json.Unmarshal(record.Body.Bytes(), &bodyResponse))
			assert.Equal(t, expectedTopStatisticsAccessedResponse, bodyResponse)
		})
		t.Run("ok with no limit", func(t *testing.T) {
			// Given
			router := NewBuilder(domain.EnvTest).WithGetTopStatisticsHandler(mockCmd(&statistics.StatisticTypeAccessed, 0, topAccessedStatistics, nil)).router
			u, err := url.Parse(fmt.Sprintf("%s/statistics/accessed", pathPrefixV1))
			require.NoError(t, err)

			// When
			record := httptest.NewRecorder()
			req := httptest.NewRequest("GET", u.String(), nil)
			router.ServeHTTP(record, req)

			// Then
			assert.Equal(t, http.StatusOK, record.Code)
			bodyResponse := GetTopStatisticsResponse{}
			require.NoError(t, json.Unmarshal(record.Body.Bytes(), &bodyResponse))
			assert.Equal(t, expectedTopStatisticsAccessedResponse, bodyResponse)
		})
		t.Run("bad request", func(t *testing.T) {
			// Given
			router := NewBuilder(domain.EnvTest).WithGetTopStatisticsHandler(mockCmd(nil, 0, topAccessedStatistics, nil)).router
			u, err := url.Parse(fmt.Sprintf("%s/statistics/accessed?limit=not-an-integer", pathPrefixV1))
			require.NoError(t, err)

			// When
			record := httptest.NewRecorder()
			req := httptest.NewRequest("GET", u.String(), nil)
			router.ServeHTTP(record, req)

			// Then
			assert.Equal(t, http.StatusBadRequest, record.Code)
		})
		t.Run("internal server error", func(t *testing.T) {
			// Given
			router := NewBuilder(domain.EnvTest).WithGetTopStatisticsHandler(mockCmd(nil, 0, topAccessedStatistics, assert.AnError)).router
			u, err := url.Parse(fmt.Sprintf("%s/statistics/accessed", pathPrefixV1))
			require.NoError(t, err)

			// When
			record := httptest.NewRecorder()
			req := httptest.NewRequest("GET", u.String(), nil)
			router.ServeHTTP(record, req)

			// Then
			assert.Equal(t, http.StatusInternalServerError, record.Code)
		})
	})
	t.Run("for shortened", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			// Given
			router := NewBuilder(domain.EnvTest).WithGetTopStatisticsHandler(mockCmd(&statistics.StatisticTypeShortened, limit, topShortenedStatistics, nil)).router
			u, err := url.Parse(fmt.Sprintf("%s/statistics/shortened?limit=%d", pathPrefixV1, limit))
			require.NoError(t, err)

			// When
			record := httptest.NewRecorder()
			req := httptest.NewRequest("GET", u.String(), nil)
			router.ServeHTTP(record, req)

			// Then
			assert.Equal(t, http.StatusOK, record.Code)
			bodyResponse := GetTopStatisticsResponse{}
			require.NoError(t, json.Unmarshal(record.Body.Bytes(), &bodyResponse))
			assert.Equal(t, expectedTopStatisticsShortenedResponse, bodyResponse)
		})
		t.Run("ok with no limit", func(t *testing.T) {
			// Given
			router := NewBuilder(domain.EnvTest).WithGetTopStatisticsHandler(mockCmd(&statistics.StatisticTypeShortened, 0, topShortenedStatistics, nil)).router
			u, err := url.Parse(fmt.Sprintf("%s/statistics/shortened", pathPrefixV1))
			require.NoError(t, err)

			// When
			record := httptest.NewRecorder()
			req := httptest.NewRequest("GET", u.String(), nil)
			router.ServeHTTP(record, req)

			// Then
			assert.Equal(t, http.StatusOK, record.Code)
			bodyResponse := GetTopStatisticsResponse{}
			require.NoError(t, json.Unmarshal(record.Body.Bytes(), &bodyResponse))
			assert.Equal(t, expectedTopStatisticsShortenedResponse, bodyResponse)
		})
		t.Run("bad request", func(t *testing.T) {
			// Given
			router := NewBuilder(domain.EnvTest).WithGetTopStatisticsHandler(mockCmd(nil, 0, topShortenedStatistics, nil)).router
			u, err := url.Parse(fmt.Sprintf("%s/statistics/shortened?limit=not-an-integer", pathPrefixV1))
			require.NoError(t, err)

			// When
			record := httptest.NewRecorder()
			req := httptest.NewRequest("GET", u.String(), nil)
			router.ServeHTTP(record, req)

			// Then
			assert.Equal(t, http.StatusBadRequest, record.Code)
		})
		t.Run("internal server error", func(t *testing.T) {
			// Given
			router := NewBuilder(domain.EnvTest).WithGetTopStatisticsHandler(mockCmd(nil, 0, topShortenedStatistics, assert.AnError)).router
			u, err := url.Parse(fmt.Sprintf("%s/statistics/shortened", pathPrefixV1))
			require.NoError(t, err)

			// When
			record := httptest.NewRecorder()
			req := httptest.NewRequest("GET", u.String(), nil)
			router.ServeHTTP(record, req)

			// Then
			assert.Equal(t, http.StatusInternalServerError, record.Code)
		})
	})
}
