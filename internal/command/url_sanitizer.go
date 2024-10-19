package command

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	// ErrInvalidURL is the error when a URL is invalid
	ErrInvalidURL error = errors.New("url is invalid")
)

// URLSanitizerCmd represents an URL sanitizer function signature
type URLSanitizerCmd func(rawURL string) (string, error)

// sanitizeURL sanitizes the given URL by trimming spaces, converting to lowercase, and normalizing.
func sanitizeURL() URLSanitizerCmd {
	return func(rawURL string) (string, error) {
		// Trim leading and trailing spaces
		cleanedURL := strings.TrimSpace(rawURL)

		// Parse the URL
		parsedURL, err := url.Parse(cleanedURL)
		if err != nil {
			return "", fmt.Errorf("%w: %w", ErrInvalidURL, err)
		}

		// Lowercase the scheme and host
		parsedURL.Scheme = strings.ToLower(parsedURL.Scheme)
		parsedURL.Host = strings.ToLower(parsedURL.Host)

		// Normalize the path
		parsedURL.Path = strings.TrimRight(parsedURL.Path, "/")
		if parsedURL.Path == "" {
			parsedURL.Path = "/"
		}

		return parsedURL.String(), nil
	}
}

// URLSanitizerCmdBuilder builds an URL sanitizercommand
func URLSanitizerCmdBuilder() URLSanitizerCmd {
	return sanitizeURL()
}
