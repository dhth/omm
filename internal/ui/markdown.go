package ui

import (
	_ "embed"

	"github.com/charmbracelet/glamour"
	"github.com/muesli/termenv"
)

var (
	//go:embed assets/gruvbox.json
	glamourJsonBytes []byte
)

func getMarkDownRenderer(wrap int) (*glamour.TermRenderer, error) {
	return glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes(glamourJsonBytes),
		glamour.WithColorProfile(termenv.TrueColor),
		glamour.WithWordWrap(wrap),
	)
}
