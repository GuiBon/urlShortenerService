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
		urlsDeleted := 10
		shortURLMock := shorturl.NewMock(t)
		shortURLMock.On("DeleteExpired", mock.Anything, timeToExpire).Return(urlsDeleted, nil)
		cmd := DeleteExpiredURLsCmdBuilder(timeToExpire, shortURLMock)

		// When
		urlsDeletedResult, err := cmd(context.Background())

		// Then
		require.NoError(t, err)
		assert.Equal(t, urlsDeleted, urlsDeletedResult)
	})
	t.Run("deletion failed", func(t *testing.T) {
		// Given
		timeToExpire := 1 * time.Hour
		shortURLMock := shorturl.NewMock(t)
		shortURLMock.On("DeleteExpired", mock.Anything, timeToExpire).Return(0, assert.AnError)
		cmd := DeleteExpiredURLsCmdBuilder(timeToExpire, shortURLMock)

		// When
		urlsDeletedResult, err := cmd(context.Background())

		// Then
		require.ErrorIs(t, err, assert.AnError)
		assert.Zero(t, urlsDeletedResult)
	})
}
