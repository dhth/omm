package theme

import (
	"errors"
	"fmt"
	"strings"
)

const (
	DefaultThemeName = themeNameGruvboxDark
)

var ErrInvalidThemeName = errors.New("invalid theme name provided")

var themes = []Theme{
	catppuccinMocha(),
	dracula(),
	githubDark(),
	gruvboxDark(),
	gruvboxDarkHard(),
	gruvboxLight(),
	monokaiClassic(),
	oneDark(),
	rosePineMoon(),
	tokyonight(),
	xcodeDark(),
}

type Theme struct {
	Name              string
	Accent1           string
	Accent2           string
	Accent3           string
	Accent4           string
	Accent5           string
	Accent6           string
	Success           string
	Danger            string
	Muted             string
	Foreground        string
	Background        string
	CategoricalColors []string
}

func All() []string {
	names := make([]string, 0, len(themes))
	for _, thm := range themes {
		names = append(names, thm.Name)
	}

	return names
}

func Get(name string) (Theme, error) {
	trimmed := strings.TrimSpace(name)

	for _, thm := range themes {
		if thm.Name == trimmed {
			return cloneTheme(thm), nil
		}
	}

	return Theme{}, fmt.Errorf("%w: %q", ErrInvalidThemeName, name)
}

func NextTheme(name string) (Theme, error) {
	return themeByOffset(name, 1)
}

func PreviousTheme(name string) (Theme, error) {
	return themeByOffset(name, -1)
}

func themeByOffset(name string, offset int) (Theme, error) {
	trimmed := strings.TrimSpace(name)
	currentIndex := -1
	for i, thm := range themes {
		if thm.Name == trimmed {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return Theme{}, fmt.Errorf("%w: %q", ErrInvalidThemeName, name)
	}

	targetIndex := (currentIndex + offset + len(themes)) % len(themes)

	return cloneTheme(themes[targetIndex]), nil
}

func cloneTheme(thm Theme) Theme {
	cp := thm
	cp.CategoricalColors = append([]string(nil), thm.CategoricalColors...)

	return cp
}
