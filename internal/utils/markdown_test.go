package utils

import (
	"encoding/json"
	"testing"

	"github.com/charmbracelet/glamour"
	"github.com/stretchr/testify/assert"
)

func TestGetGlamourStyleFromFile(t *testing.T) {
	gotOption := glamour.WithStylesFromJSONBytes(glamourJsonBytes)
	renderer, err := glamour.NewTermRenderer(gotOption)
	assert.NoError(t, err)
	assert.NotNil(t, renderer)

	_, err = renderer.Render("a")
	assert.NoError(t, err)
}

func TestGlamourStylesFileIsValid(t *testing.T) {
	got := json.Valid(glamourJsonBytes)
	assert.True(t, got)
}
