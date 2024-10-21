package http

import (
	"urlShortenerService/domain"
	"urlShortenerService/internal/usecase"

	"github.com/gin-gonic/gin"
)

const pathPrefixV1 = "/api/url-shortener/v1"

// Builder holds the gin Engine
type Builder struct {
	router *gin.Engine
}

// NewBuilder creates a Builder
func NewBuilder(env domain.Environment) *Builder {
	switch env {
	case domain.EnvTest:
		gin.SetMode(gin.TestMode)
	case domain.EnvProduction:
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	return &Builder{
		router: gin.Default(),
	}
}

// BuildRouter builds the gin Engine router
func (b *Builder) BuildRouter(createShortenURLCmd usecase.CreateShortenURLCmd, getOriginalURLCmd usecase.GetOriginalURLCmd,
	getStatisticsForURLCmd usecase.GetStatisticsForURLCmd) *gin.Engine {
	return b.
		WithSwaggerHandler().
		WithV1HealthHandler().
		WithV1CreateShortenURLHandler(createShortenURLCmd).
		WithGetOriginalURLHandler(getOriginalURLCmd).
		WithGetStatisticsForURLHandler(getStatisticsForURLCmd).
		router
}
