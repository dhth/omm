package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mvdan.cc/xurls/v2"
)

func TestExtractURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		// success
		{
			name:  "just a url",
			input: `https://someurl.com`,
			expected: []string{
				"https://someurl.com",
			},
		},
		{
			name:  "a url with path",
			input: `https://someurl.com/path/1`,
			expected: []string{
				"https://someurl.com/path/1",
			},
		},
		{
			name:  "a url with query parameters",
			input: `https://someurl.com?param=value`,
			expected: []string{
				"https://someurl.com?param=value",
			},
		},
		{
			name: "two urls",
			input: `https://someurl.com
https://anotherurl.com`,
			expected: []string{
				"https://someurl.com",
				"https://anotherurl.com",
			},
		},
		{
			name: "urls in a paragraph",
			input: `A paragraph full of details, containing urls like https://someurl.com/path?query=value
and https://anotherurl.com/path?query=value at several points.`,
			expected: []string{
				"https://someurl.com/path?query=value",
				"https://anotherurl.com/path?query=value",
			},
		},
		{
			name: "urls ending with commas and braces",
			input: `A paragraph full of details, containing urls
(eg. https://someurl.com/path?query=value, https://anotherurl.com/path?query=value)
at several points.`,
			expected: []string{
				"https://someurl.com/path?query=value",
				"https://anotherurl.com/path?query=value",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			rxStrict := xurls.Strict()

			got := ExtractURLs(rxStrict, tt.input)

			assert.Equal(t, tt.expected, got)
		})
	}
}
