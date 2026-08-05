package theme

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Contrast floors are calibrated to each color's UI role and the established
	// visual hierarchy of omm's intentionally adapted themes.
	minimumForegroundContrast  = 4.5
	minimumProminentContrast   = 4.0
	minimumCategoricalContrast = 3.5
	minimumMutedContrast       = 3.0
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
			expectedName: "github-light",
		},
		{
			name:         "next theme wraps around",
			currentTheme: "xcode-dark",
			expectedName: "catppuccin-latte",
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
			expectedName: "github-light",
		},
		{
			name:         "previous theme wraps around",
			currentTheme: "catppuccin-latte",
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

func TestColorsHaveSufficientContrast(t *testing.T) {
	for _, themeName := range All() {
		t.Run(themeName, func(t *testing.T) {
			thm, err := Get(themeName)
			require.NoError(t, err)

			foregroundContrast := assertMinimumContrast(t, "foreground", thm.Foreground, thm.Background, minimumForegroundContrast)
			mutedContrast := assertMinimumContrast(t, "muted", thm.Muted, thm.Background, minimumMutedContrast)
			assert.Lessf(t, mutedContrast, foregroundContrast,
				"muted %s should have lower contrast than foreground %s against background %s",
				thm.Muted, thm.Foreground, thm.Background,
			)

			prominentColors := []struct {
				name  string
				color string
			}{
				{name: "accent1", color: thm.Accent1},
				{name: "accent2", color: thm.Accent2},
				{name: "accent3", color: thm.Accent3},
				{name: "accent4", color: thm.Accent4},
				{name: "accent5", color: thm.Accent5},
				{name: "accent6", color: thm.Accent6},
				{name: "success", color: thm.Success},
				{name: "danger", color: thm.Danger},
			}
			for _, color := range prominentColors {
				assertMinimumContrast(t, color.name, color.color, thm.Background, minimumProminentContrast)
			}

			require.NotEmpty(t, thm.CategoricalColors)
			for i, color := range thm.CategoricalColors {
				assertMinimumContrast(t, fmt.Sprintf("categorical color %d", i), color, thm.Background, minimumCategoricalContrast)
			}
		})
	}
}

func assertMinimumContrast(t *testing.T, name, foreground, background string, minimum float64) float64 {
	t.Helper()

	foregroundLuminance, err := relativeLuminance(foreground)
	require.NoErrorf(t, err, "%s has invalid color %q", name, foreground)
	backgroundLuminance, err := relativeLuminance(background)
	require.NoErrorf(t, err, "background is an invalid color %q", background)

	// WCAG contrast ranges from 1:1 for identical luminance to 21:1 for
	// black against white. The lighter color is always the numerator.
	// See https://www.w3.org/WAI/WCAG21/Understanding/contrast-minimum.html
	lighter := math.Max(foregroundLuminance, backgroundLuminance)
	darker := math.Min(foregroundLuminance, backgroundLuminance)
	ratio := (lighter + 0.05) / (darker + 0.05)

	assert.GreaterOrEqualf(t, ratio, minimum,
		"%s %s against background %s has contrast %.2f:1; want at least %.1f:1",
		name, foreground, background, ratio, minimum,
	)

	return ratio
}

func relativeLuminance(hexColor string) (float64, error) {
	if len(hexColor) != 7 || hexColor[0] != '#' {
		return 0, fmt.Errorf("expected #RRGGBB color")
	}

	// Hex colors contain gamma-encoded sRGB channels. Convert them to linear
	// light before calculating luminance. The channel weights below approximate
	// human brightness perception, which is most sensitive to green and least to blue.
	channels := make([]float64, 3)
	for i := range channels {
		value, err := strconv.ParseUint(hexColor[1+i*2:3+i*2], 16, 8)
		if err != nil {
			return 0, fmt.Errorf("couldn't parse RGB channel: %w", err)
		}

		channel := float64(value) / 255
		if channel <= 0.04045 {
			channels[i] = channel / 12.92
		} else {
			channels[i] = math.Pow((channel+0.055)/1.055, 2.4)
		}
	}

	return 0.2126*channels[0] + 0.7152*channels[1] + 0.0722*channels[2], nil
}
