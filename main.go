package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/dhth/omm/cmd"
)

const (
	author    = "@dhth"
	issuesURL = "https://github.com/dhth/omm/issues"
)

var version = "dev"

func main() {
	v := version
	if version == "dev" {
		info, ok := debug.ReadBuildInfo()
		if ok {
			v = info.Main.Version
		}
	}
	err := cmd.Execute(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		followUp, isUnexpected := cmd.GetFollowUp(err)

		if len(followUp) > 0 {
			fmt.Fprintf(os.Stderr, `
%s
`, followUp)
		}

		if isUnexpected {
			fmt.Fprintf(os.Stderr, `
---

This isn't supposed to happen; let %s know about this error via %s.
`, author, issuesURL)
		}
		os.Exit(1)
	}
}
