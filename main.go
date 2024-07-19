package main

import (
	"github.com/dhth/omm/cmd"
	"runtime/debug"
)

var (
	version = "dev"
)

func main() {
	v := version
	if version == "dev" {
		info, ok := debug.ReadBuildInfo()
		if ok {
			v = info.Main.Version
		}
	}
	cmd.Execute(v)
}
