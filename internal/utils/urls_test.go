package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractURIs(t *testing.T) {
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
		{
			name: "uris with custom schemes",
			input: `
slack link: slack://open?team=T12345678&id=C12345678
obsidian link without space after the colon:obsidian://open?vault=VAULT&file=FILE
maps: maps://?q=Central+Park,New+York
a scheme from the future: skynet://someresource?query=param
`,
			expected: []string{
				"slack://open?team=T12345678&id=C12345678",
				"obsidian://open?vault=VAULT&file=FILE",
				"maps://?q=Central+Park,New+York",
				"skynet://someresource?query=param",
			},
		},
		{
			name: "known schemes that use `:`",
			input: `
mail: mailto:example@example.com
telephone: tel:+1234567890
spotify: spotify:track:6rqhFgbbKwnb9MLmUQDhG6
facetime: facetime:example@example.com
facetime-audio: facetime-audio:example@example.com
`,
			expected: []string{
				"mailto:example@example.com",
				"tel:+1234567890",
				"spotify:track:6rqhFgbbKwnb9MLmUQDhG6",
				"facetime:example@example.com",
				"facetime-audio:example@example.com",
			},
		},
		// failures
		{
			name:  "doesn't match a uri without a scheme",
			input: `someurl.com`,
		},
		{
			name:  "unknown schemes that use `:`",
			input: "unknown: unknown:example@example.com",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			rgx := GetURIRegex()

			got := ExtractURIs(rgx, tt.input)

			assert.Equal(t, tt.expected, got)
			if len(tt.expected) > 0 {
				assert.Equal(t, tt.expected, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}
