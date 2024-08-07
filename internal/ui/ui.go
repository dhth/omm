package ui

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func RenderUI(db *sql.DB, config Config) {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Fprintf(os.Stderr, "fatal error: %s", err.Error())
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(InitialModel(db, config), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Something went wrong %s", err)
	}
}
