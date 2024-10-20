package command

import (
	"crypto/sha1"

	"github.com/jxskiss/base62"
)

// SlugGeneratorCmd represents a slug generator function signature
type SlugGeneratorCmd func(url string) string

// generateSlug generates a consistent slug for a given URL
func generateSlug(slugLenght int) SlugGeneratorCmd {
	return func(url string) string {
		// Generate SHA-1 hash of the URL
		hasher := sha1.New()
		hasher.Write([]byte(url))
		hashBytes := hasher.Sum(nil)

		// Encode the hash bytes using Base62
		base62Hash := base62.EncodeToString(hashBytes)

		// Shorten the slug to only 8 characters
		slug := base62Hash
		if len(slug) > slugLenght {
			slug = slug[:slugLenght]
		}

		return slug
	}
}

// SlugGeneratorCmdBuilder builds a slug generator command
func SlugGeneratorCmdBuilder(slugLenght int) SlugGeneratorCmd {
	return generateSlug(slugLenght)
}
