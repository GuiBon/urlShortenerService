package http

import (
	"fmt"
	"net/http"
	"net/url"
	"urlShortenerService/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// GetStatisticsForURLResponse holds the JSON body response structure
type GetStatisticsForURLResponse struct {
	URL              string `json:"url"`
	ShortenedCounter int    `json:"shortened_counter"`
	AccessedCounter  int    `json:"accessed_counter"`
}

// WithGetStatisticsForURLHandler register the get statistics for URL API in the router of the HTTP builder
func (b *Builder) WithGetStatisticsForURLHandler(cmd usecase.GetStatisticsForURLCmd) *Builder {
	b.router.GET(fmt.Sprintf("%s/statistics", pathPrefixV1), getStatisticsForURLHandler(cmd))
	return b
}

// getStatisticsForURLHandler retrieves statistics for a given URL
func getStatisticsForURLHandler(cmd usecase.GetStatisticsForURLCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		var encodedURL string
		encodedURL, encodedURLExists := c.GetQuery("encoded_url")
		if !encodedURLExists {
			c.JSON(http.StatusBadRequest, CreateAPIError(ApiError{
				Name:        "bad_request",
				Description: "no URL given as query parameter",
				Hint:        "add an URL in query parameter name 'encoded_url'",
			}, nil))
			return
		}

		url, err := url.QueryUnescape(encodedURL)
		if err != nil {
			c.JSON(http.StatusBadRequest, CreateAPIError(ApiError{
				Name:        "bad_request",
				Description: "unable to unescape given URL",
				Hint:        "badly encoded URL",
			}, err))
			return
		}

		statistics, err := cmd(c.Request.Context(), url)
		switch err {
		case nil:
			c.JSON(http.StatusOK, GetStatisticsForURLResponse{
				URL:              statistics.URL,
				ShortenedCounter: statistics.ShortenedCounter,
				AccessedCounter:  statistics.AccessedCounter,
			})
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
