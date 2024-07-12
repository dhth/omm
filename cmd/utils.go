package cmd

import (
	"os"
	"os/user"
	"strings"
)

const (
	fallbackEditor = "vi"
)

func expandTilde(path string) string {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			os.Exit(1)
		}
		return strings.Replace(path, "~", usr.HomeDir, 1)
	}
	return path
}

func getUserConfiguredEditor() string {
	editor := os.Getenv("EDITOR")
	if editor != "" {
		return editor
	}

	editor = os.Getenv("VISUAL")
	if editor != "" {
		return editor
	}

	return fallbackEditor
}
