package http

import (
	"fmt"
	"net/http"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// CreateShortenURLRequest holds the JSON body request structure
type CreateShortenURLRequest struct {
	OriginalURL string `json:"original_url" binding:"required"`
}

// CreateShortenURLResponse holds the JSON body response structure
type CreateShortenURLResponse struct {
	ShortURL string `json:"short_url"`
}

// WithV1CreateShortenURLHandler register the create shorten URL API in the router of the HTTP builder
func (b *Builder) WithV1CreateShortenURLHandler(cmd usecase.CreateShortenURLCmd) *Builder {
	b.router.POST(fmt.Sprintf("%s/shorten", pathPrefixV1), v1CreateShortenURLHandler(cmd))
	return b
}

// v1CreateShortenURLHandler creates a shorten URL of the given one
func v1CreateShortenURLHandler(cmd usecase.CreateShortenURLCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		var createShortenURLRequest CreateShortenURLRequest
		err := c.ShouldBindJSON(&createShortenURLRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, CreateAPIError(ApiError{
				Name:        "bad_request",
				Description: "can't parse JSON body",
				Hint:        "the body should be JSON with application/json and required fields",
			}, err))
			return
		}

		shortenedURL, err := cmd(c.Request.Context(), createShortenURLRequest.OriginalURL)
		switch err {
		case nil:
			c.JSON(http.StatusOK, CreateShortenURLResponse{ShortURL: shortenedURL})
			return
		case command.ErrInvalidURL:
			c.JSON(http.StatusUnprocessableEntity, CreateAPIError(ApiError{
				Name:        "unprocessable_entity",
				Description: "the given original_url is invalid",
				Hint:        "the URL should respect the RFC: https://datatracker.ietf.org/doc/html/rfc1738 ",
			}, err))
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
