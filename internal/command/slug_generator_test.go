package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSlug(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		// Given
		url := "https://example.com"
		cmd := SlugGeneratorCmdBuilder(8)

		// When
		slug := cmd(url)

		// Then
		assert.Len(t, slug, 8)
	})
	t.Run("slug generation is consistent", func(t *testing.T) {
		// Given
		url := "https://example.com"
		slugs := map[string]interface{}{}
		cmd := SlugGeneratorCmdBuilder(8)

		// When
		for i := 0; i < 5; i++ {
			slug := cmd(url)
			slugs[slug] = true
		}

		// Then
		assert.Len(t, slugs, 1)

	})
}
