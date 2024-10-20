package http

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateAPIError(t *testing.T) {
	t.Run("release mode", func(t *testing.T) {
		// Given
		gin.SetMode(gin.ReleaseMode)
		apiError := ApiError{
			Name:        "test name",
			Description: "test description",
			Hint:        "test hint",
		}

		// When
		fullAPIError := CreateAPIError(apiError, assert.AnError)

		// Then
		assert.Equal(t, apiError.Name, fullAPIError.Name)
		assert.Equal(t, apiError.Description, fullAPIError.Description)
		assert.Equal(t, apiError.Hint, fullAPIError.Hint)
		assert.Nil(t, fullAPIError.Debug)
	})
	t.Run("non release mode", func(t *testing.T) {
		// Given
		gin.SetMode(gin.DebugMode)
		apiError := ApiError{
			Name:        "test name",
			Description: "test description",
			Hint:        "test hint",
		}

		// When
		fullAPIError := CreateAPIError(apiError, assert.AnError)

		// Then
		assert.Equal(t, apiError.Name, fullAPIError.Name)
		assert.Equal(t, apiError.Description, fullAPIError.Description)
		assert.Equal(t, apiError.Hint, fullAPIError.Hint)
		assert.Equal(t, assert.AnError, fullAPIError.Debug)

	})
}
