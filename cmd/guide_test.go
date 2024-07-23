package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContext(t *testing.T) {
	for _, entry := range guideEntries {
		got, err := getContext(entry.summary)
		assert.Nil(t, err)
		assert.NotEmpty(t, got)
	}
}
