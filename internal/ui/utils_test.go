package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSummaryWithNewPrefix(t *testing.T) {
	testCases := []struct {
		name      string
		summary   string
		newPrefix string
		expected  string
	}{
		{
			name:      "empty summary",
			summary:   "",
			newPrefix: "new",
			expected:  "new: ",
		},
		{
			name:      "just whitespace",
			summary:   "  ",
			newPrefix: "new",
			expected:  "new: ",
		},
		{
			name:      "just a colon",
			summary:   ":",
			newPrefix: "new",
			expected:  "new: ",
		},
		{
			name:      "just summary content",
			summary:   "this is a task",
			newPrefix: "new",
			expected:  "new: this is a task",
		},
		{
			name:      "just a prefix",
			summary:   "old:",
			newPrefix: "new",
			expected:  "new: ",
		},
		{
			name:      "prefix and summary content",
			summary:   "old: this is a task",
			newPrefix: "new",
			expected:  "new: this is a task",
		},
		{
			name:      "prefix and summary content without space",
			summary:   "old:this is a task",
			newPrefix: "new",
			expected:  "new:this is a task",
		},
		{
			name:      "summary with two colons",
			summary:   "old: this: is a task",
			newPrefix: "new",
			expected:  "new: this: is a task",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := getSummaryWithNewPrefix(tt.summary, tt.newPrefix)

			assert.Equal(t, tt.expected, got)
		})
	}
}
