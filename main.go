package main

import (
	"homelab-inventory/cmd"
	"homelab-inventory/internal/logging"
	"homelab-inventory/internal/version"
)

var (
	Version   = "1.0.0"
	Commit    = "unset"
	BuildTime = "unset"
	GoVersion = "unset"
)

func main() {
	version.Set(Version, Commit, BuildTime, GoVersion)
	logging.Init(true)
	cmd.Execute()
}
