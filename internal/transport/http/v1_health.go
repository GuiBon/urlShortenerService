package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WithV1HealthHandler register the health API in the router of the HTTP builder
func (b *Builder) WithV1HealthHandler() *Builder {
	b.router.GET(fmt.Sprintf("%s/health", pathPrefixV1), v1HealthHandler())
	return b
}

// v1HealthHandler informs about health status of the service
func v1HealthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}
