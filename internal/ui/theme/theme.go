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
	monokaiClassic(),
	oneDark(),
	rosePineMoon(),
	tokyonight(),
	xcodeDark(),
}

type Theme struct {
	Name         string
	Primary      string
	Secondary    string
	Tertiary     string
	Quaternary   string
	Quinary      string
	Success      string
	Error        string
	Muted        string
	Text         string
	Background   string
	PrefixColors []string
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
			return thm, nil
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

	return themes[targetIndex], nil
}
