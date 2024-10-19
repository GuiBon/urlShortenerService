package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitizeURL(t *testing.T) {
	scenarios := []struct {
		Name           string
		RawURL         string
		ExpectedOutput string
		ExpectError    bool
	}{
		{Name: "nominal", RawURL: "https://example.com/path", ExpectedOutput: "https://example.com/path"},
		{Name: "nominal with trailing slash", RawURL: "https://example.com/path/", ExpectedOutput: "https://example.com/path"},
		{Name: "nominal with spaces", RawURL: " https://example.com/path ", ExpectedOutput: "https://example.com/path"},
		{Name: "nominal with no path", RawURL: "https://example.com", ExpectedOutput: "https://example.com/"},
		{Name: "nominal with uppercase", RawURL: "HTTPS://EXAMPLE.COM/PATH", ExpectedOutput: "https://example.com/PATH"},
		{Name: "error while parsing url", RawURL: "://example.com", ExpectedOutput: "", ExpectError: true},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			// Given
			cmd := URLSanitizerCmdBuilder()

			// When
			sanitizedURL, err := cmd(scenario.RawURL)

			// Then
			require.Equal(t, err != nil, scenario.ExpectError)
			assert.Equal(t, scenario.ExpectedOutput, sanitizedURL)
		})
	}
}
