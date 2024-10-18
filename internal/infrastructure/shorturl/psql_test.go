package shorturl

import (
	"os"
	"testing"
	"urlShortenerService/internal/infrastructure/config"

	"github.com/stretchr/testify/require"
)

func TestPSQLStore(t *testing.T) {
	os.Setenv("env", "test")
	defer os.Unsetenv("env")
	conf, err := config.Load()
	require.NoError(t, err)
	store, err := NewPSQLStore(conf.Database)
	require.NoError(t, err)

	RunStoreTests(t, store)
}
