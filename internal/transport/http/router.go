package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewRouter creates a router
func NewRouter() *gin.Engine {
	router := gin.Default()

	// Endpoint GET /health to ensure that the service is up
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return router
}
