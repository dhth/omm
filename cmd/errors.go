package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dhth/omm/internal/ui/theme"
)

func GetFollowUp(err error) (string, bool) {
	var followUp string
	var isUnexpected bool

	switch {
	case errors.Is(err, errCouldntGetHomeDir):
		isUnexpected = true
	case errors.Is(err, errCouldntSetupGuide):
		isUnexpected = true
	case errors.Is(err, theme.ErrInvalidThemeName):
		followUp = fmt.Sprintf("Tip: valid themes are [%s]", strings.Join(theme.All(), ", "))
	}

	return followUp, isUnexpected
}
