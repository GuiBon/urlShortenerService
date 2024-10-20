package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// WithSwaggerHandler register the swagger API in the router of the HTTP builder
func (b *Builder) WithSwaggerHandler() *Builder {
	b.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.yaml")))
	b.router.GET("/docs/swagger.yaml", func(c *gin.Context) {
		c.File("/app/docs/swagger.yaml")
	})
	return b
}
