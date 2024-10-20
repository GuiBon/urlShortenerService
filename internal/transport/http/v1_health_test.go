package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"urlShortenerService/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithV1HealthHandler(t *testing.T) {
	// Given
	router := NewBuilder(domain.EnvTest).WithV1HealthHandler().router

	// When
	u, err := url.Parse(fmt.Sprintf("%s/health", pathPrefixV1))
	require.NoError(t, err)
	record := httptest.NewRecorder()
	req := httptest.NewRequest("GET", u.String(), nil)
	router.ServeHTTP(record, req)

	// Then
	assert.Equal(t, http.StatusOK, record.Code)
}
