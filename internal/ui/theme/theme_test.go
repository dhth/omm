package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultThemeIsValid(t *testing.T) {
	// GIVEN
	// WHEN
	_, err := Get(DefaultThemeName)

	// THEN
	assert.NoError(t, err)
}

func TestNextThemeWorksForAllThemes(t *testing.T) {
	for _, themeName := range All() {
		// GIVEN
		// WHEN
		_, err := NextTheme(themeName)

		// THEN
		require.NoError(t, err)
	}
}

func TestPreviousThemeWorksForAllThemes(t *testing.T) {
	for _, themeName := range All() {
		// GIVEN
		// WHEN
		_, err := PreviousTheme(themeName)

		// THEN
		require.NoError(t, err)
	}
}

func TestNextTheme(t *testing.T) {
	testCases := []struct {
		name         string
		currentTheme string
		expectedName string
		expectedErr  error
	}{
		{
			name:         "next theme in middle of list",
			currentTheme: "github-dark",
			expectedName: "gruvbox-dark",
		},
		{
			name:         "next theme wraps around",
			currentTheme: "xcode-dark",
			expectedName: "catppuccin-mocha",
		},
		{
			name:         "next theme trims whitespace",
			currentTheme: "  dracula  ",
			expectedName: "github-dark",
		},
		{
			name:         "next theme fails for unknown input",
			currentTheme: "does-not-exist",
			expectedErr:  ErrInvalidThemeName,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			currentTheme := tt.currentTheme

			// WHEN
			nextTheme, err := NextTheme(currentTheme)

			// THEN
			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedName, nextTheme.Name)
		})
	}
}

func TestPreviousTheme(t *testing.T) {
	testCases := []struct {
		name         string
		currentTheme string
		expectedName string
		expectedErr  error
	}{
		{
			name:         "previous theme in middle of list",
			currentTheme: "gruvbox-dark",
			expectedName: "github-dark",
		},
		{
			name:         "previous theme wraps around",
			currentTheme: "catppuccin-mocha",
			expectedName: "xcode-dark",
		},
		{
			name:         "previous theme trims whitespace",
			currentTheme: "  github-dark  ",
			expectedName: "dracula",
		},
		{
			name:         "previous theme fails for unknown input",
			currentTheme: "unknown-theme",
			expectedErr:  ErrInvalidThemeName,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			currentTheme := tt.currentTheme

			// WHEN
			previousTheme, err := PreviousTheme(currentTheme)

			// THEN
			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedName, previousTheme.Name)
		})
	}
}
