package command

import (
	"errors"
	"unicode"
)

var (
	// ErrInvalidSlugLenght is the error when a slug lenght is invalid
	ErrInvalidSlugLenght error = errors.New("slug lenght is invalid")
	// ErrInvalidSlugNonAlphanumeric is the error when a slug is invalid because it has non alphanumeric character
	ErrInvalidSlugNonAlphanumeric error = errors.New("slug is invalid because of non alphanumeric character")
)

// SlugValidatorCmd represents a slug validator function signature
type SlugValidatorCmd func(slug string) error

// validateSlug ensures that a slug is valid
func validateSlug(slugLenght int) SlugValidatorCmd {
	return func(slug string) error {
		if len(slug) > slugLenght {
			return ErrInvalidSlugLenght
		}
		for _, char := range slug {
			if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
				return ErrInvalidSlugNonAlphanumeric
			}
		}
		return nil
	}
}

// SlugValidatorCmdBuilder builds a slug validator command
func SlugValidatorCmdBuilder(slugLenght int) SlugValidatorCmd {
	return validateSlug(slugLenght)
}
