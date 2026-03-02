package ui

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
	"github.com/dhth/omm/internal/ui/theme"
	"github.com/muesli/termenv"
)

func getMarkDownRenderer(thm theme.Theme, wrap int) (*glamour.TermRenderer, error) {
	return glamour.NewTermRenderer(
		glamour.WithStyles(glamourStyleConfig(thm)),
		glamour.WithColorProfile(termenv.TrueColor),
		glamour.WithPreservedNewLines(),
		glamour.WithWordWrap(wrap),
	)
}

func glamourStyleConfig(thm theme.Theme) ansi.StyleConfig {
	boolPtr := func(v bool) *bool { return &v }
	stringPtr := func(s string) *string { return &s }
	uintPtr := func(v uint) *uint { return &v }

	return ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "",
				BlockSuffix: "",
				Color:       stringPtr(thm.Text),
			},
			Margin: uintPtr(2),
		},
		BlockQuote: ansi.StyleBlock{
			Indent:      uintPtr(1),
			IndentToken: stringPtr("┃ "),
		},
		Paragraph: ansi.StyleBlock{},
		List: ansi.StyleList{
			LevelIndent: 2,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       stringPtr(thm.Primary),
				Bold:        boolPtr(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "# ",
				Color:  stringPtr(thm.Primary),
				Bold:   boolPtr(true),
			},
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "## ",
				Color:  stringPtr(thm.Tertiary),
			},
		},
		H3: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "### ",
				Color:  stringPtr(thm.Tertiary),
			},
		},
		H4: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "#### ",
				Color:  stringPtr(thm.Tertiary),
			},
		},
		H5: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: stringPtr(thm.Secondary),
			},
		},
		H6: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: stringPtr(thm.Success),
				Bold:  boolPtr(false),
			},
		},
		Text:          ansi.StylePrimitive{},
		Strikethrough: ansi.StylePrimitive{CrossedOut: boolPtr(true)},
		Emph: ansi.StylePrimitive{
			Color:  stringPtr(thm.Tertiary),
			Italic: boolPtr(true),
		},
		Strong: ansi.StylePrimitive{
			Color: stringPtr(thm.Primary),
			Bold:  boolPtr(true),
		},
		HorizontalRule: ansi.StylePrimitive{
			Color:  stringPtr(thm.Muted),
			Format: "\n--------\n",
		},
		Item:        ansi.StylePrimitive{BlockPrefix: "• "},
		Enumeration: ansi.StylePrimitive{BlockPrefix: ". "},
		Task: ansi.StyleTask{
			Ticked:   "[✔] ",
			Unticked: "[ ] ",
		},
		Link: ansi.StylePrimitive{
			Color: stringPtr(thm.Tertiary),
		},
		LinkText: ansi.StylePrimitive{
			Color: stringPtr(thm.Secondary),
			Bold:  boolPtr(true),
		},
		Image: ansi.StylePrimitive{
			Color:     stringPtr("132"),
			Underline: boolPtr(true),
		},
		ImageText: ansi.StylePrimitive{
			Color:  stringPtr("245"),
			Format: "Image: {{.text}} →",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: stringPtr(thm.Quaternary),
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			Theme: chromaTheme(thm.Name),
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color: stringPtr(thm.Text),
				},
			},
		},
		Table: ansi.StyleTable{
			CenterSeparator: stringPtr("┼"),
			ColumnSeparator: stringPtr("│"),
			RowSeparator:    stringPtr("─"),
		},
		DefinitionList: ansi.StyleBlock{},
		DefinitionTerm: ansi.StylePrimitive{},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n🠶 ",
		},
		HTMLBlock: ansi.StyleBlock{},
		HTMLSpan:  ansi.StyleBlock{},
	}
}

func chromaTheme(themeName string) string {
	switch themeName {
	case "catppuccin-mocha":
		return "catppuccin-mocha"
	case "dracula":
		return "dracula"
	case "github-dark":
		return "github-dark"
	case "gruvbox-dark":
		return "gruvbox"
	case "monokai-classic":
		return "monokai"
	case "onedark":
		return "onedark"
	case "rose-pine-moon":
		return "rose-pine-moon"
	case "tokyonight":
		return "tokyonight-night"
	case "xcode-dark":
		return "xcode-dark"
	default:
		return "gruvbox"
	}
}
