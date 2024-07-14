package cmd

import (
	"errors"
	"os"
	"os/user"
	"strings"

	"github.com/dhth/omm/internal/ui"
)

var (
	showContextEnvVarMisconfiguredErr = errors.New("OMM_SHOW_CONTEXT can only be one of these values: 0/1/true/false")
	listDensityEnvVarMisconfiguredErr = errors.New("OMM_LIST_DENSITY can only be one of these values: compact/spacious")
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

func getUserConfiguredShowContext(defaultVal bool) (bool, error) {
	sc := os.Getenv("OMM_SHOW_CONTEXT")
	if sc == "" {
		return defaultVal, nil
	}

	if sc == "0" || sc == "false" {
		return false, nil
	}
	if sc == "1" || sc == "true" {
		return true, nil
	}

	return false, showContextEnvVarMisconfiguredErr
}

func getUserConfiguredListDensity(defaultVal string) (ui.ListDensityType, error) {
	ld := os.Getenv("OMM_LIST_DENSITY")
	if ld == "" {
		ld = defaultVal
	}

	switch ld {
	case ui.CompactDensityVal:
		return ui.Compact, nil
	case ui.SpaciousDensityVal:
		return ui.Spacious, nil
	}

	return ui.Compact, listDensityEnvVarMisconfiguredErr
}
