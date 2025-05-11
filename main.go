package main

import (
	"homelab-inventory/cmd"
	"homelab-inventory/internal/logging"
)

func main() {
	logging.Init()
	cmd.Execute()
}
