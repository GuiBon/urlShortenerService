package http

import (
	"strings"
	"urlShortenerService/internal/usecase"

	"github.com/gin-gonic/gin"
)

// WithGetOriginalURLForceHandler register the force get original URL API in the router of the HTTP builder
func (b *Builder) WithGetOriginalURLForceHandler(cmd usecase.GetOriginalURLCmd) *Builder {
	b.router.GET("/:slug/force", cleanForceURLPath(), getOriginalURLHandler(cmd))
	return b
}

// cleanForceURLPath cleans the path before retrieving an original URL given a slug
func cleanForceURLPath() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cleanedPath string
		if strings.HasSuffix(c.Request.URL.Path, "/force") {
			cleanedPath = strings.TrimSuffix(c.Request.URL.Path, "/force")
		} else { // We assume that there is query parameters
			// Few chances out of several billions but a slug might be /forceAAA
			// So we need to ensure that there is query parameter after before replacing
			// Otherwise http://localhost:8080/forceAAA/force?redirect=true will be replaced by http://localhost:8080AAA/force?redirect=true
			cleanedPath = strings.Replace(c.Request.URL.Path, "/force?", "?", 1)
		}
		c.Request.URL.Path = cleanedPath
		c.Next()
	}
}
