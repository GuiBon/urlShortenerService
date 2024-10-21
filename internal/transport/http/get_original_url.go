package http

import (
	"net/http"
	"strconv"
	"urlShortenerService/internal/command"
	"urlShortenerService/internal/infrastructure/malwarescanner"
	"urlShortenerService/internal/infrastructure/shorturl"
	"urlShortenerService/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// GetOriginalURLResponse holds the JSON body response structure
type GetOriginalURLResponse struct {
	OriginalURL string `json:"original_url"`
}

// WithGetOriginalURLHandler register the get original URL API in the router of the HTTP builder
func (b *Builder) WithGetOriginalURLHandler(cmd usecase.GetOriginalURLCmd) *Builder {
	b.router.GET("/:slug", getOriginalURLHandler(cmd))
	return b
}

// getOriginalURLHandler retrieves an original URL given a slug
func getOriginalURLHandler(cmd usecase.GetOriginalURLCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		var redirect bool
		redirectStr, redirectQueryParamsExists := c.GetQuery("redirect")
		if redirectQueryParamsExists {
			var err error
			redirect, err = strconv.ParseBool(redirectStr)
			if err != nil {
				glog.Warningf("failed to parse the 'redirect' query parameter value, redirection ignored")
			}
		}

		originalURL, err := cmd(c.Request.Context(), slug)
		switch err {
		case nil:
			if redirect {
				c.Redirect(http.StatusFound, originalURL)
			} else {
				c.JSON(http.StatusOK, GetOriginalURLResponse{OriginalURL: originalURL})
			}
			return
		case malwarescanner.ErrMalswareURL:
			c.JSON(http.StatusForbidden, CreateAPIError(ApiError{
				Name:        "forbidden",
				Description: "a malware has been detected within the URL",
				Hint:        "if you really want to continue use /force API",
			}, err))
			return
		case shorturl.ErrNotFound:
			c.JSON(http.StatusNotFound, CreateAPIError(ApiError{
				Name:        "not_found",
				Description: "no URL found associated to the given slug",
				Hint:        "the slug might be incorrect or expired",
			}, err))
			return
		case command.ErrInvalidSlugLenght, command.ErrInvalidSlugNonAlphanumeric:
			c.JSON(http.StatusUnprocessableEntity, CreateAPIError(ApiError{
				Name:        "unprocessable_entity",
				Description: "the given slug is invalid",
				Hint:        "the slug should be alpha numeric and less than the configuration setted maximal lenght",
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
