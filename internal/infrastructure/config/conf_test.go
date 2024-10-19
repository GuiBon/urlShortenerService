package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		// Given
		os.Setenv("env", "test")
		defer os.Unsetenv("env")

		// When
		conf, err := Load()

		// Then
		require.NoError(t, err)
		assert.NotEmpty(t, conf)
	})
	t.Run("config file not fount", func(t *testing.T) {
		// Given
		os.Setenv("env", "unknown")
		defer os.Unsetenv("env")

		// When
		conf, err := Load()

		// Then
		require.Error(t, err)
		assert.Empty(t, conf)
	})
}

func TestToConnString(t *testing.T) {
	// Given
	conf := PSQLConnConfig{
		User:     "user",
		Password: "password",
		Host:     "localhost",
		Port:     5432,
		DbName:   "dbtest",
	}

	// When
	connString := conf.ToConnString()

	// Then
	assert.Equal(t, "postgres://user:password@localhost:5432/dbtest?sslmode=disable", connString)
}

func TestCreateBaseURL(t *testing.T) {
	t.Run("with port", func(t *testing.T) {
		// Given
		cfg := ServerDomainConfig{
			Scheme: "http",
			Domain: "localhost",
			Port:   8080,
		}

		// When
		baseURL := cfg.CreateBaseURL()

		// Then
		assert.Equal(t, "http://localhost:8080", baseURL)
	})
	t.Run("without port", func(t *testing.T) {
		// Given
		cfg := ServerDomainConfig{
			Scheme: "https",
			Domain: "www.example.com",
		}

		// When
		baseURL := cfg.CreateBaseURL()

		// Then
		assert.Equal(t, "https://www.example.com", baseURL)
	})
}
