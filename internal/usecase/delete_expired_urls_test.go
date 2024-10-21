package usecase

import (
	"context"
	"testing"
	"time"
	"urlShortenerService/internal/infrastructure/shorturl"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteExpiredURLsCmdBuilder(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		// Given
		timeToExpire := 1 * time.Hour
		slugsDeleted := []string{"2zv8a2Im", "1eJSWjFM", "UsIJeS1D", "K11q8dTj", "Sd7k2eDU"}
		shortURLMock := shorturl.NewMock(t)
		shortURLMock.On("DeleteExpired", mock.Anything, timeToExpire).Return(slugsDeleted, nil)
		cmd := DeleteExpiredURLsCmdBuilder(timeToExpire, shortURLMock)

		// When
		slugsDeletedResult, err := cmd(context.Background())

		// Then
		require.NoError(t, err)
		assert.Equal(t, slugsDeleted, slugsDeletedResult)
	})
	t.Run("deletion failed", func(t *testing.T) {
		// Given
		timeToExpire := 1 * time.Hour
		shortURLMock := shorturl.NewMock(t)
		shortURLMock.On("DeleteExpired", mock.Anything, timeToExpire).Return([]string{}, assert.AnError)
		cmd := DeleteExpiredURLsCmdBuilder(timeToExpire, shortURLMock)

		// When
		slugsDeletedResult, err := cmd(context.Background())

		// Then
		require.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, slugsDeletedResult)
	})
}
