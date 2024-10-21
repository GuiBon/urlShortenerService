package http

import (
	"fmt"
	"net/http"
	"strconv"
	"urlShortenerService/internal/infrastructure/statistics"
	"urlShortenerService/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// GetTopStatisticsResponse holds the JSON body response structure
type GetTopStatisticsResponse struct {
	URLs []getTopStatisticsForURLResponse `json:"urls"`
}

type getTopStatisticsForURLResponse struct {
	URL              string `json:"url"`
	ShortenedCounter int    `json:"shortened_counter,omitempty"`
	AccessedCounter  int    `json:"accessed_counter,omitempty"`
}

// WithGetTopStatisticsHandler register the get top statistics API in the router of the HTTP builder
func (b *Builder) WithGetTopStatisticsHandler(cmd usecase.GetTopStatisticsCmd) *Builder {
	b.router.GET(fmt.Sprintf("%s/statistics/accessed", pathPrefixV1), getTopStatisticsHandler(statistics.StatisticTypeAccessed, cmd))
	b.router.GET(fmt.Sprintf("%s/statistics/shortened", pathPrefixV1), getTopStatisticsHandler(statistics.StatisticTypeShortened, cmd))
	return b
}

// getTopStatisticsHandler retrieves top statistics
func getTopStatisticsHandler(statType statistics.StatisticType, cmd usecase.GetTopStatisticsCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		var resultLimit int
		resultLimitStr, resultLimitExists := c.GetQuery("limit")
		if resultLimitExists {
			var err error
			resultLimit, err = strconv.Atoi(resultLimitStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, CreateAPIError(ApiError{
					Name:        "bad_request",
					Description: "invalid query parameter 'limit'",
					Hint:        "'limit' value is not an integer",
				}, err))
				return
			}
		}

		topStatistics, err := cmd(c.Request.Context(), statType, int64(resultLimit))
		switch err {
		case nil:
			var response = GetTopStatisticsResponse{}
			for _, topStatistic := range topStatistics {
				response.URLs = append(response.URLs, getTopStatisticsForURLResponse{
					URL:              topStatistic.URL,
					AccessedCounter:  topStatistic.AccessedCounter,
					ShortenedCounter: topStatistic.ShortenedCounter,
				})
			}
			c.JSON(http.StatusOK, response)
			return
		default:
			glog.Error(err)
			c.JSON(http.StatusInternalServerError, CreateAPIError(ApiError{
				Name:        "internal_server_error",
				Description: "unknown error",
				Hint:        "if you are the application owner, please check the logs for more details",
			}, err))
			return
		}
	}
}
