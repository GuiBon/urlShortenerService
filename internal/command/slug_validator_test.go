package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatorSlug(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		// Given
		slug := "zTw34enA"
		cmd := SlugValidatorCmdBuilder(8)

		// When
		err := cmd(slug)

		// Then
		assert.NoError(t, err)
	})
	t.Run("with a short slug", func(t *testing.T) {
		// Given
		slug := "z"
		cmd := SlugValidatorCmdBuilder(8)

		// When
		err := cmd(slug)

		// Then
		assert.NoError(t, err)
	})
	t.Run("slug is invalid", func(t *testing.T) {
		t.Run("because of lenght", func(t *testing.T) {
			// Given
			slug := "zTw34enAh"
			cmd := SlugValidatorCmdBuilder(8)

			// When
			err := cmd(slug)

			// Then
			assert.ErrorIs(t, err, ErrInvalidSlugLenght)
		})
		t.Run("because of non alphanumeric character", func(t *testing.T) {
			// Given
			slug := "zTw+4enA"
			cmd := SlugValidatorCmdBuilder(8)

			// When
			err := cmd(slug)

			// Then
			assert.ErrorIs(t, err, ErrInvalidSlugNonAlphanumeric)
		})
	})
}
