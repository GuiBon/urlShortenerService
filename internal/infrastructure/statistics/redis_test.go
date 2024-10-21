package statistics

import (
	"strconv"
	"testing"
	"urlShortenerService/internal/infrastructure/config"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/require"
)

func TestRedisStore(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	port, err := strconv.Atoi(mr.Port())
	require.NoError(t, err)
	store, err := NewRedisStore(config.RedisConfig{Host: mr.Host(), Port: port, MaxResults: 10})
	require.NoError(t, err)

	RunStoreTests(t, store)
}
