package utils

import (
	_ "embed"

	"github.com/charmbracelet/glamour"
	"github.com/muesli/termenv"
)

//go:embed assets/gruvbox.json
var glamourJSONBytes []byte

func GetMarkDownRenderer(wrap int) (*glamour.TermRenderer, error) {
	return glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes(glamourJSONBytes),
		glamour.WithColorProfile(termenv.TrueColor),
		glamour.WithWordWrap(wrap),
	)
}
