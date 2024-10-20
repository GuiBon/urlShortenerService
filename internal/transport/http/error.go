package http

import "github.com/gin-gonic/gin"

// ApiError represents a minimal HTTP error
type ApiError struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Hint        string `json:"hint"`
}

// fullAPIError represents a full HTTP error with debug
type fullAPIError struct {
	ApiError
	Debug error `json:"debug"`
}

// CreateAPIError creates an API error with debug value for non release mode only
func CreateAPIError(apiError ApiError, err error) fullAPIError {
	if gin.Mode() == gin.ReleaseMode {
		return fullAPIError{ApiError: apiError}
	}
	return fullAPIError{ApiError: apiError, Debug: err}
}
