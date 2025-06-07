package cmd

import (
	"os"
	"os/user"
	"strings"
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

func getUserConfiguredEditor(defaultVal string) string {
	editor := os.Getenv("OMM_EDITOR")
	if editor != "" {
		return editor
	}

	editor = os.Getenv("EDITOR")
	if editor != "" {
		return editor
	}

	editor = os.Getenv("VISUAL")
	if editor != "" {
		return editor
	}

	return defaultVal
}
