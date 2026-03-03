package ui

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/dhth/omm/internal/ui/theme"
)

func RenderUI(db *sql.DB, config Config, thm theme.Theme) {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Fprintf(os.Stderr, "fatal error: %s", err.Error())
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(InitialModel(db, config, thm))
	if _, err := p.Run(); err != nil {
		log.Fatalf("Something went wrong %s", err)
	}
}
