package ui

import (
	"github.com/charmbracelet/glamour"
	"github.com/muesli/termenv"
)

func getMarkDownRenderer(wrap int) (*glamour.TermRenderer, error) {

	var margin uint = 2
	dracula := glamour.DraculaStyleConfig
	dracula.Document.BlockPrefix = ""
	dracula.H1.Prefix = ""
	dracula.H2.Prefix = ""
	dracula.H3.Prefix = ""
	dracula.H4.Prefix = ""
	dracula.Document.Margin = &margin

	return glamour.NewTermRenderer(
		glamour.WithColorProfile(termenv.TrueColor),
		glamour.WithStyles(dracula),
		glamour.WithWordWrap(wrap),
	)
}
