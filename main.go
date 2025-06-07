package main

import (
	"os"
	"runtime/debug"

	"github.com/dhth/omm/internal/cmd"
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
		os.Exit(1)
	}
}
