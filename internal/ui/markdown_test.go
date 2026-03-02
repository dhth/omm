package ui

import (
	"testing"

	chromastyles "github.com/alecthomas/chroma/v2/styles"
	"github.com/dhth/omm/internal/ui/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMarkDownRenderer(t *testing.T) {
	for _, name := range theme.All() {
		t.Run(name, func(t *testing.T) {
			// GIVEN
			thm, err := theme.Get(name)
			require.NoError(t, err)

			// WHEN
			renderer, err := getMarkDownRenderer(thm, 80)

			// THEN
			require.NoError(t, err)
			require.NotNil(t, renderer)

			_, err = renderer.Render("# Some markdown text")
			assert.NoError(t, err)
		})
	}
}

func TestAllMappedChromaThemesAreAvailable(t *testing.T) {
	themes := theme.All()
	themes = append(themes, "unknown")

	for _, name := range themes {
		t.Run(name, func(t *testing.T) {
			// GIVEN
			mapped := chromaTheme(name)

			// WHEN
			_, ok := chromastyles.Registry[mapped]

			// THEN
			assert.Truef(t, ok, "mapped chroma theme %q for theme %q not found", mapped, name)
		})
	}
}
